package audit

import (
	"context"
	"net/netip"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v3/transport"
	"github.com/google/uuid"
	"github.com/zhixingdx/public/api/apicommon/v1"
	"github.com/zhixingdx/public/api/config"
	"github.com/zhixingdx/public/pkg/middleware/meta"
	"go.opentelemetry.io/otel/trace"
)

// DeploymentSource 使用稳定部署配置标识来源；同部署扩缩容、进程重启不会改变幂等命名空间。
func DeploymentSource(c *config.Config, service string) string {
	name, env := c.GetServer().GetName(), c.GetServer().GetEnv()
	if name == "" {
		name = service
	}
	if env == "" {
		env = "default"
	}
	return service + "/" + env + ":" + name
}

// NewMetadataRecord 在交给异步出口前固定事件身份、发生时间与链路信息，重试复用该信封。
func NewMetadataRecord(ctx context.Context, source string, call *apicommon.BizCall, category string) *Record {
	id := uuid.NewString()
	record := &Record{EventID: id, OperationID: id, SourceDeploymentID: source, BizCall: call,
		OccurredAt: time.Now().UnixMilli(), PolicyVersion: 1, SchemaVersion: 1, Category: category, CaptureMode: "metadata"}
	clientID, _ := meta.GetClientID(ctx)
	channel, _ := meta.GetChannel(ctx)
	platform, _ := meta.GetPlatform(ctx)
	record.ClientID = optionalContextText(clientID, 127)
	record.Channel = optionalContextText(channel, 63)
	record.Platform = optionalContextText(platform, 63)
	if realm, err := meta.GetRealm(ctx); err == nil {
		record.Realm = optionalContextText(realm, 63)
	}
	// 终端头只从 HTTP 请求冻结；gRPC 的 User-Agent 描述传输库，不能冒充业务客户端。
	if tr, ok := transport.FromServerContext(ctx); ok && tr.Kind() == transport.KindHTTP {
		record.UserAgent = optionalContextText(meta.GetUserAgent(ctx), 1024)
		ip := meta.GetRequestHeaderVal(ctx, "X-Real-IP")
		if ip == "" {
			ip = meta.GetRequestHeaderVal(ctx, "X-Forwarded-For")
		}
		// IP 是请求报告的上下文；不拆分 XFF 代理链，也不能据此授予权限。
		if address, err := netip.ParseAddr(ip); err == nil && address.Zone() == "" {
			record.ClientIP = optionalContextText(address.String(), 45)
		}
	}
	span := trace.SpanContextFromContext(ctx)
	if span.IsValid() {
		record.TraceID, record.SpanID = span.TraceID().String(), span.SpanID().String()
	}
	return record
}

// SetSourceMetadata 在异步入队前固定生产者本机环境与版本，不从来源标识或远端请求反推。
func SetSourceMetadata(record *Record, environment, version string) {
	if record == nil {
		return
	}
	record.SourceEnvironment = optionalContextText(environment, 63)
	record.ServiceVersion = optionalContextText(version, 127)
}

// optionalContextText 丢弃超长、非法编码和控制字符的自报上下文，避免可选信息影响主业务。
func optionalContextText(value string, limit int) *string {
	if !utf8.ValidString(value) || len(value) > limit || strings.IndexFunc(value, unicode.IsControl) >= 0 {
		return nil
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
