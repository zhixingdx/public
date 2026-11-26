package server

import (
	"context"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/metadata"
	"github.com/go-kratos/kratos/v3/transport"
	kgrpc "github.com/go-kratos/kratos/v3/transport/grpc"
	common "github.com/zhixingdx/oipivot-api/apicommon/v1"
	iam "github.com/zhixingdx/oipivot-api/base/iam/v1"
	"github.com/zhixingdx/oipivot-api/config"
	"github.com/zhixingdx/oipivot-api/fault"
	"github.com/zhixingdx/oipivot/backend/pkg/lang"
	"github.com/zhixingdx/oipivot/backend/pkg/middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/structpb"
)

const trustedRuntimeInputLimit = 128 << 10
const trustedRuntimeOutputLimit = 2 << 20
const trustedRuntimeTimeout = 10 * time.Second

// TrustedRuntimePolicy 复用所属服务 HTTP 入口的认证配置及权限列表，不建立独立操作授权目录。
type TrustedRuntimePolicy struct {
	Auth *config.Auth
	List *auth.ListConf
}

// trustedRuntimeServer 持有当前实例的方法快照和固定本机连接，不接受跨服务地址。
type trustedRuntimeServer struct {
	common.UnimplementedTrustedRuntimeServer
	methods     map[string]protoreflect.MethodDescriptor
	permissions map[string]middleware.Middleware
	conn        *grpc.ClientConn
	lang        *lang.CommonLang
}

// RegisterTrustedRuntime 在业务 RPC 注册完成后增加受信授权入口，返回固定本机连接的清理函数。
// 旧 RPC 保持原系统调用契约；仅本入口执行严格员工授权，网络来源仍由既有受信内网边界保证。
func RegisterTrustedRuntime(srv *kgrpc.Server, listener *config.ServerItem, policies map[string]TrustedRuntimePolicy, l *lang.CommonLang) (func(), error) {
	if srv == nil || listener == nil || l == nil {
		return nil, fmt.Errorf("受信调用入口缺少服务、监听配置或语言依赖")
	}
	if _, exists := srv.GetServiceInfo()["api.apicommon.TrustedRuntime"]; exists {
		return nil, fmt.Errorf("受信调用入口已注册")
	}
	endpoint, err := trustedRuntimeLoopback(listener)
	if err != nil {
		return nil, err
	}
	runtime := &trustedRuntimeServer{methods: map[string]protoreflect.MethodDescriptor{}, permissions: map[string]middleware.Middleware{}, lang: l}
	connections := map[string]*grpc.ClientConn{}
	registered := false
	cleanup := func() {
		for _, conn := range connections {
			_ = conn.Close()
		}
	}
	defer func() {
		if !registered {
			cleanup()
		}
	}()
	for serviceName, policy := range policies {
		if policy.Auth == nil || policy.Auth.Type != auth.TypeIam || strings.TrimSpace(policy.Auth.Addr) == "" {
			return nil, fmt.Errorf("受信服务 %s 必须使用已有 IAM 认证配置", serviceName)
		}
		service, ok := srv.GetServiceInfo()[serviceName]
		if !ok || serviceName == "api.apicommon.TrustedRuntime" {
			return nil, fmt.Errorf("受信服务 %s 未注册或为递归入口", serviceName)
		}
		descriptor, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName(serviceName))
		if err != nil {
			return nil, fmt.Errorf("受信服务 %s 缺少本地契约", serviceName)
		}
		desc, ok := descriptor.(protoreflect.ServiceDescriptor)
		if !ok {
			return nil, fmt.Errorf("受信服务 %s 契约类型非法", serviceName)
		}
		iamConn := connections[policy.Auth.Addr]
		if iamConn == nil {
			iamConn, err = kgrpc.NewClient(context.Background(), kgrpc.WithEndpoint(policy.Auth.Addr), kgrpc.WithTimeout(6*time.Second), kgrpc.WithMiddleware(metadata.Client()))
			if err != nil {
				return nil, fmt.Errorf("创建现有 IAM 授权连接失败: %w", err)
			}
			connections[policy.Auth.Addr] = iamConn
		}
		deps := auth.TrustedStaffDependencies{GetStaffCache: iam.NewBaseIamManageClient(iamConn).GetStaffCache, ResolveTenantRealm: iam.NewBaseIamManagePlusClient(iamConn).ResolveTenantRealm, CheckStaffPermission: iam.NewBaseIamPortalPlusClient(iamConn).CheckStaffPermission}
		check := auth.NewTrustedStaffMiddleware(deps, policy.List, l)
		for _, method := range service.Methods {
			if method.IsClientStream || method.IsServerStream {
				continue
			}
			md := desc.Methods().ByName(protoreflect.Name(method.Name))
			if md == nil || md.IsStreamingClient() || md.IsStreamingServer() {
				return nil, fmt.Errorf("受信服务 %s 方法契约不一致", serviceName)
			}
			operation := "/" + serviceName + "/" + method.Name
			runtime.methods[operation] = md
			runtime.permissions[operation] = check
		}
	}
	if len(runtime.methods) == 0 {
		return nil, fmt.Errorf("受信调用入口没有可用方法")
	}
	conn, err := kgrpc.NewClient(context.Background(), kgrpc.WithEndpoint(endpoint), kgrpc.WithTimeout(trustedRuntimeTimeout), kgrpc.WithMiddleware(metadata.Client()), kgrpc.WithOptions(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(trustedRuntimeOutputLimit), grpc.MaxCallSendMsgSize(trustedRuntimeInputLimit))))
	if err != nil {
		return nil, fmt.Errorf("创建受信本机连接失败: %w", err)
	}
	runtime.conn = conn
	common.RegisterTrustedRuntimeServer(srv, runtime)
	registered = true
	return func() { _ = conn.Close(); cleanup() }, nil
}

// trustedRuntimeLoopback 仅把本实例的 TCP 通配或本机监听转换为固定回环地址，不推导远端服务地址。
func trustedRuntimeLoopback(listener *config.ServerItem) (string, error) {
	if listener.Network != "" && listener.Network != "tcp" {
		return "", fmt.Errorf("受信调用入口仅支持 TCP 监听")
	}
	host, port, err := net.SplitHostPort(listener.Addr)
	if err != nil {
		return "", fmt.Errorf("受信调用监听地址非法")
	}
	number, err := strconv.Atoi(port)
	if err != nil || number < 1 || number > 65535 {
		return "", fmt.Errorf("受信调用监听端口非法")
	}
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	} else if host == "::" {
		host = "::1"
	} else if host == "localhost" {
		host = "127.0.0.1"
	} else if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
		return "", fmt.Errorf("受信调用入口必须绑定回环或通配监听")
	}
	return net.JoinHostPort(host, port), nil
}

// trustedRuntimeTransport 仅覆盖实际目标操作，保持原认证头与传输属性供统一权限中间件读取。
type trustedRuntimeTransport struct {
	transport.Transporter
	operation string
}

// Operation 返回已按当前实例注册表验证的目标操作。
func (t trustedRuntimeTransport) Operation() string { return t.operation }

// Invoke 先校验当前实例方法与资源上界，再按目标操作授权，最后经原 RPC 完整执行验证和业务逻辑。
func (s *trustedRuntimeServer) Invoke(ctx context.Context, req *common.TrustedRuntimeRequest) (*common.TrustedRuntimeReply, error) {
	bad := func() error { return fault.ErrorBadRequest("%s", s.lang.GetParameterErrorMsg(ctx)) }
	if req == nil || len(req.Operation) > 511 || req.Input == nil || proto.Size(req.Input) > trustedRuntimeInputLimit || !trustedRuntimeSafeNumbers(req.Input.AsMap(), 0) {
		return nil, bad()
	}
	method, ok := s.methods[req.Operation]
	if !ok {
		return nil, fault.ErrorForbidden("%s", s.lang.GetNoPermissionMsg(ctx))
	}
	check, ok := s.permissions[req.Operation]
	if !ok || check == nil {
		return nil, fault.ErrorForbidden("%s", s.lang.GetNoPermissionMsg(ctx))
	}
	original, ok := transport.FromServerContext(ctx)
	if !ok {
		return nil, bad()
	}
	ctx = transport.NewServerContext(ctx, trustedRuntimeTransport{Transporter: original, operation: req.Operation})
	ctx, cancel := context.WithTimeout(ctx, trustedRuntimeTimeout)
	defer cancel()
	result, err := check(func(callCtx context.Context, _ any) (any, error) {
		data, err := protojson.Marshal(req.Input)
		if err != nil || len(data) > trustedRuntimeInputLimit {
			return nil, bad()
		}
		input := dynamicpb.NewMessage(method.Input())
		if err = protojson.Unmarshal(data, input); err != nil {
			return nil, bad()
		}
		output := dynamicpb.NewMessage(method.Output())
		if err = s.conn.Invoke(callCtx, req.Operation, input, output); err != nil {
			// 连接和基础设施原文可能包含地址或内部实现；业务授权/校验错误及取消仍保持原协议。
			switch status.Code(err) {
			case codes.Unavailable, codes.Internal, codes.Unknown, codes.ResourceExhausted:
				return nil, fault.ErrorServiceUnavailable("%s", s.lang.GetRequestErrorMsg(ctx))
			}
			return nil, err
		}
		return s.encodeOutput(ctx, output)
	})(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*common.TrustedRuntimeReply), nil
}

// trustedRuntimeSafeNumbers 拒绝 Struct 中会丢失整数精度的数字及过深嵌套，64 位 protobuf 整数由 protojson 字符串保留。
func trustedRuntimeSafeNumbers(value any, depth int) bool {
	if depth > 64 {
		return false
	}
	switch current := value.(type) {
	case float64:
		return !math.IsNaN(current) && !math.IsInf(current, 0) && (math.Trunc(current) != current || math.Abs(current) <= 9007199254740991)
	case map[string]any:
		for _, child := range current {
			if !trustedRuntimeSafeNumbers(child, depth+1) {
				return false
			}
		}
	case []any:
		for _, child := range current {
			if !trustedRuntimeSafeNumbers(child, depth+1) {
				return false
			}
		}
	}
	return true
}

// encodeOutput 保持契约 JSONName 与缺省值语义，拒绝超界响应及 Struct 数字精度丢失。
func (s *trustedRuntimeServer) encodeOutput(ctx context.Context, output proto.Message) (*common.TrustedRuntimeReply, error) {
	failed := func() error { return fault.ErrorServiceUnavailable("%s", s.lang.GetRequestErrorMsg(ctx)) }
	if proto.Size(output) > trustedRuntimeOutputLimit {
		return nil, failed()
	}
	data, err := protojson.MarshalOptions{EmitDefaultValues: true}.Marshal(output)
	if err != nil || len(data) > trustedRuntimeOutputLimit {
		return nil, failed()
	}
	result := &structpb.Struct{}
	if err = protojson.Unmarshal(data, result); err != nil || !trustedRuntimeSafeNumbers(result.AsMap(), 0) {
		return nil, failed()
	}
	return &common.TrustedRuntimeReply{Output: result}, nil
}
