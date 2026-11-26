package audit

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	public "github.com/zhixingdx/oipivot-api/base/audit/v1"
	"github.com/zhixingdx/oipivot-api/config"
	"github.com/zhixingdx/oipivot/backend/pkg/middleware/meta"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"
)

// ProducerConfig 可声明无租户列的归属；启用证书模式时同时配置服务端信任根和租户证书。
type ProducerConfig struct {
	UnscopedTenant string                      `json:"unscopedTenant,omitempty"`
	CAFile         string                      `json:"caFile,omitempty"`
	ServerName     string                      `json:"serverName,omitempty"`
	Tenants        map[string]TenantCredential `json:"tenants,omitempty"`
}

// TenantCredential 用挂载文件提供一把租户专用的工作负载身份。
type TenantCredential struct {
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
}

// ProducerClient 使用部署配置中的固定端点；可信内网共用连接，证书模式按租户隔离连接。
type ProducerClient struct {
	unscopedTenant     string
	networkClient      public.BaseAuditPortalPlusClient
	clients            map[string]public.BaseAuditPortalPlusClient
	connections        []*grpc.ClientConn
	SourceDeploymentID string
	SourceEnvironment  string
	ServiceVersion     string
}

// producerSettings 逐字段校验动态配置，只有明确未启用 TLS 时才允许可信内网连接。
func producerSettings(extra *structpb.Struct) (settings ProducerConfig, hasTLS bool, err error) {
	fields := extra.GetFields()
	for key := range fields {
		if strings.EqualFold(key, "auditProducer") && key != "auditProducer" {
			return settings, false, fmt.Errorf("未知配置 %q，应为 auditProducer", key)
		}
	}
	raw, exists := fields["auditProducer"]
	if !exists {
		return settings, false, nil
	}
	if raw.GetStructValue() == nil {
		return settings, false, fmt.Errorf("auditProducer 须为对象")
	}
	for key, value := range raw.GetStructValue().GetFields() {
		switch key {
		case "unscopedTenant":
			settings.UnscopedTenant = value.GetStringValue()
			if !meta.ValidTenantID(settings.UnscopedTenant) {
				return settings, false, fmt.Errorf("无效的无租户列实体归属")
			}
		case "caFile":
			hasTLS = true
			settings.CAFile = value.GetStringValue()
			if settings.CAFile == "" {
				return settings, false, fmt.Errorf("auditProducer.caFile 须为非空字符串")
			}
		case "serverName":
			hasTLS = true
			settings.ServerName = value.GetStringValue()
			if settings.ServerName == "" {
				return settings, false, fmt.Errorf("auditProducer.serverName 须为非空字符串")
			}
		case "tenants":
			hasTLS = true
			settings.Tenants, err = producerCredentials(value.GetStructValue())
			if err != nil {
				return settings, false, err
			}
		default:
			return settings, false, fmt.Errorf("auditProducer 的未知字段 %q", key)
		}
	}
	return settings, hasTLS, nil
}

// producerCredentials 只接受显式租户与完整证书路径，拒绝大小写或类型不符的配置。
func producerCredentials(entries *structpb.Struct) (map[string]TenantCredential, error) {
	if len(entries.GetFields()) == 0 {
		return nil, fmt.Errorf("auditProducer.tenants 须为非空租户对象")
	}
	tenants := make(map[string]TenantCredential, len(entries.Fields))
	for tenant, entry := range entries.Fields {
		if !meta.ValidTenantID(tenant) {
			return nil, fmt.Errorf("无效的证书租户 %q", tenant)
		}
		fields := entry.GetStructValue().GetFields()
		cert := TenantCredential{CertFile: fields["certFile"].GetStringValue(), KeyFile: fields["keyFile"].GetStringValue()}
		if len(fields) != 2 || cert.CertFile == "" || cert.KeyFile == "" {
			return nil, fmt.Errorf("租户 %q 的证书配置只接受非空字符串 certFile 和 keyFile", tenant)
		}
		tenants[tenant] = cert
	}
	return tenants, nil
}

// NewProducerClient 校验所有身份后才交付客户端，部分租户配置错误时收回已建连接。
func NewProducerClient(c *config.Config, service string) (*ProducerClient, error) {
	endpoint := c.GetServices()["base-audit"]
	if endpoint == "" || service == "" {
		return nil, fmt.Errorf("%w: 缺少审计端点或来源服务", ErrUnavailable)
	}
	settings, hasTLS, err := producerSettings(c.GetExtra())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	client := &ProducerClient{
		unscopedTenant: settings.UnscopedTenant, SourceDeploymentID: DeploymentSource(c, service),
		SourceEnvironment: c.GetServer().GetEnv(), ServiceVersion: c.GetServer().GetVersion(),
	}
	if !hasTLS {
		conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		client.connections = []*grpc.ClientConn{conn}
		client.networkClient = public.NewBaseAuditPortalPlusClient(conn)
		return client, nil
	}
	if settings.CAFile == "" || settings.ServerName == "" || len(settings.Tenants) == 0 {
		return nil, fmt.Errorf("%w: 须同时配置 caFile、serverName 及租户证书", ErrUnavailable)
	}
	if settings.UnscopedTenant != "" {
		if _, ok := settings.Tenants[settings.UnscopedTenant]; !ok {
			return nil, fmt.Errorf("%w: 无租户列的实体须绑定已配置的租户证书", ErrUnavailable)
		}
	}
	pem, err := os.ReadFile(settings.CAFile)
	if err != nil {
		return nil, fmt.Errorf("读取审计服务端 CA: %w", err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("%w: 无效的审计服务端 CA", ErrUnavailable)
	}
	client.clients = make(map[string]public.BaseAuditPortalPlusClient, len(settings.Tenants))
	seen := make(map[string]bool)
	for tenant, credential := range settings.Tenants {
		conn, err := dialTenant(endpoint, tenant, settings.ServerName, roots, credential, seen)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("初始化租户 %q 的审计身份: %w", tenant, err)
		}
		client.connections = append(client.connections, conn)
		client.clients[tenant] = public.NewBaseAuditPortalPlusClient(conn)
	}
	return client, nil
}

// dialTenant 校验工作负载证书，再创建仅使用该身份的连接池。
func dialTenant(endpoint, tenant, serverName string, roots *x509.CertPool, credential TenantCredential, seen map[string]bool) (*grpc.ClientConn, error) {
	if !meta.ValidTenantID(tenant) {
		return nil, fmt.Errorf("%w: 无效的租户标识", ErrUnavailable)
	}
	certificate, err := tls.LoadX509KeyPair(credential.CertFile, credential.KeyFile)
	if err != nil {
		return nil, err
	}
	leaf, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		return nil, err
	}
	if len(leaf.URIs) != 1 || leaf.URIs[0].Scheme != "spiffe" || leaf.URIs[0].Host == "" {
		return nil, fmt.Errorf("%w: 证书须包含唯一 SPIFFE URI", ErrUnavailable)
	}
	uri := leaf.URIs[0].String()
	if seen[uri] {
		return nil, fmt.Errorf("%w: 两个租户使用了相同的证书身份", ErrUnavailable)
	}
	seen[uri] = true
	transport := credentials.NewTLS(&tls.Config{
		MinVersion: tls.VersionTLS13, RootCAs: roots, ServerName: serverName,
		Certificates: []tls.Certificate{certificate},
	})
	return grpc.NewClient(endpoint, grpc.WithTransportCredentials(transport))
}

// route 按已冻结的实体归属选择连接，证书模式中未登记的租户不可发送。
func (c *ProducerClient) route(tenant string) public.BaseAuditPortalPlusClient {
	if c == nil || !meta.ValidTenantID(tenant) {
		return nil
	}
	if c.networkClient != nil {
		return c.networkClient
	}
	return c.clients[tenant]
}

// tenantClient 只发送本次事件的资源租户，覆盖调用链中已有的租户声明。
func (c *ProducerClient) tenantClient(ctx context.Context, tenant string) (public.BaseAuditPortalPlusClient, context.Context, error) {
	client := c.route(tenant)
	if client == nil {
		return nil, ctx, ErrUnavailable
	}
	md := metadata.Pairs(meta.TenantIDHeader, tenant)
	return client, metadata.NewOutgoingContext(ctx, md), nil
}

// Close 释放该来源持有的所有 gRPC 连接；运行时会在排空工作协程后调用。
func (c *ProducerClient) Close() {
	if c != nil {
		for _, conn := range c.connections {
			_ = conn.Close()
		}
	}
}
