// Package utils 通用工具包
package utils

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware/metadata"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	k "github.com/go-kratos/kratos/v3/transport/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
)

// GRPCKeepaliveOptions 返回 gRPC keepalive 配置选项（客户端）
// 保持连接活跃，检测连接状态，自动重连
// 适用于所有环境，建议始终启用
//
// 重要说明：
//   - Time: 30 秒是推荐的 ping 间隔，避免触发服务端的 "too_many_pings" 错误
//   - 如果服务端配置了 keepalive，可以适当降低到 20-30 秒
//   - 如果服务端未配置 keepalive，建议使用 30 秒或更长，避免超过服务端默认限制（2 小时内最多 2 次）
//
// 使用示例:
//
//	k.WithOptions(
//		grpc.WithStatsHandler(&tracing.ClientHandler{}),
//		utils.GRPCKeepaliveOptions()...,
//	)
func GRPCKeepaliveOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second, // 每 30 秒发送一次 keepalive ping（避免 too_many_pings 错误）
			Timeout:             3 * time.Second,  // keepalive ping 超时时间
			PermitWithoutStream: true,             // 即使没有活跃的流也发送 keepalive
		}),
	}
}

// GRPCServerKeepaliveOptions 返回 gRPC keepalive 配置选项（服务端）
// 允许客户端发送 keepalive ping，避免 "too_many_pings" 错误
// 服务端需要配置 keepalive 才能正确处理客户端的 ping 请求
//
// 配置说明：
//   - MaxConnectionIdle: 空闲连接的最大保持时间（5 分钟）
//   - MaxConnectionAge: 连接的最大存活时间（30 分钟）
//   - MaxConnectionAgeGrace: 连接关闭的宽限时间（5 秒）
//   - Time: 服务端 ping 客户端的时间间隔（2 小时，通常不需要）
//   - Timeout: ping 超时时间（20 秒）
//   - MinTime: 允许客户端发送 ping 的最小间隔（30 秒，与客户端配置一致）
//   - PermitWithoutStream: 允许在没有活跃流时发送 ping
//
// 使用示例:
//
//	import (
//		kratosGrpc "github.com/go-kratos/kratos/v3/transport/grpc"
//		"google.golang.org/grpc"
//		"google.golang.org/grpc/keepalive"
//	)
//
//	opts = append(opts, kratosGrpc.Options(
//		grpc.KeepaliveParams(keepalive.ServerParameters{
//			MaxConnectionIdle:     5 * time.Minute,
//			MaxConnectionAge:      30 * time.Minute,
//			MaxConnectionAgeGrace: 5 * time.Second,
//			Time:                  2 * time.Hour,
//			Timeout:               20 * time.Second,
//		}),
//		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
//			MinTime:             30 * time.Second,
//			PermitWithoutStream: true,
//		}),
//	))
func GRPCServerKeepaliveOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     5 * time.Minute,  // 空闲连接的最大保持时间
			MaxConnectionAge:      30 * time.Minute, // 连接的最大存活时间
			MaxConnectionAgeGrace: 5 * time.Second,  // 连接关闭的宽限时间
			Time:                  2 * time.Hour,    // 服务端 ping 客户端的时间间隔（通常不需要）
			Timeout:               20 * time.Second, // ping 超时时间
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             30 * time.Second, // 允许客户端每 30 秒发送一次 ping（与客户端配置一致）
			PermitWithoutStream: true,             // 允许在没有活跃流时发送 ping
		}),
	}
}

// MergeGRPCOptions 合并多个 gRPC DialOption 切片
// 用于简化客户端代码中的选项合并
//
// 使用示例:
//
//	k.WithOptions(
//		utils.MergeGRPCOptions(
//			[]grpc.DialOption{grpc.WithStatsHandler(&tracing.ClientHandler{})},
//			utils.GRPCKeepaliveOptions(),
//			utils.GRPCK8sOptions(), // 可选
//		)...,
//	)
func MergeGRPCOptions(opts ...[]grpc.DialOption) []grpc.DialOption {
	var result []grpc.DialOption
	for _, opt := range opts {
		result = append(result, opt...)
	}
	return result
}

// GRPCK8sOptions 返回 K8s 环境优化的 gRPC 配置选项
// 包含负载均衡、重试策略等，专门针对 K8s Service 优化
// 解决 K8s 环境中 DNS 解析失败和连接重试的问题
//
// DNS 解析和重试机制说明：
//  1. gRPC 使用 DNS resolver 解析域名，默认会缓存 DNS 结果
//  2. 连接失败时，gRPC 会触发 DNS resolver 重新解析（ResolveNow）
//  3. DNS resolver 会定期刷新（默认 30 分钟），但连接失败时会立即触发刷新
//  4. 使用 round_robin 负载均衡时，会尝试所有解析到的 IP 地址
//  5. 如果所有 IP 都连接失败，会等待 DNS 刷新后重试
//
// 使用示例（K8s 环境）:
//
//	k.WithTimeout(utils.GRPCK8sTimeout), // K8s 内网环境建议 10 秒超时
//	k.WithOptions(
//		grpc.WithStatsHandler(&tracing.ClientHandler{}),
//		utils.GRPCKeepaliveOptions()..., // 建议同时启用 keepalive
//		utils.GRPCK8sOptions()...,      // 添加 K8s 优化配置
//	)
func GRPCK8sOptions() []grpc.DialOption {
	return []grpc.DialOption{
		// 服务配置 - 负载均衡和重试策略
		// 适用于 K8s Service，支持 DNS 解析重试和连接重试
		// round_robin 负载均衡会尝试所有解析到的 IP 地址
		// 连接失败时会触发 DNS 重新解析，获取最新的 IP 列表
		grpc.WithDefaultServiceConfig(`{
			"loadBalancingConfig": [{"round_robin":{}}],
			"methodConfig": [{
				"name": [{"service": ""}],
				"retryPolicy": {
					"maxAttempts": 3,
					"initialBackoff": "0.1s",
					"maxBackoff": "1s",
					"backoffMultiplier": 2.0,
					"retryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED", "RESOURCE_EXHAUSTED"]
				},
				"timeout": "30s"
			}]
		}`),
		// 注意：如果不配置 timeout，gRPC 默认没有超时限制
		// 这意味着请求可能会无限期等待，直到完成或连接断开
		// 建议始终配置 timeout，避免请求挂起
	}
}

// GRPCK8sTimeout K8s 内网环境建议的 gRPC 连接建立超时时间
//
// 重要说明：这是连接建立的超时时间，不是连接保持的超时时间
//
// k.WithTimeout() 的作用：
//   - 控制从调用 NewClient() 开始，到连接成功建立的最大等待时间
//   - 如果在这个时间内无法建立连接，NewClient() 会返回错误
//   - 连接建立成功后，这个超时时间就不再起作用了
//
// 连接建立过程包括：
//  1. DNS 解析（查询 CoreDNS，将服务名解析为 IP）- 通常 < 100ms
//  2. TCP 连接建立（三次握手）- 通常 < 10ms
//  3. gRPC 握手（协议协商）- 通常 < 50ms
//
// 正常情况下的时间消耗：
//   - 内网 DNS 解析：通常 10-100ms
//   - TCP 连接建立：通常 1-10ms
//   - gRPC 握手：通常 10-50ms
//   - 总计：正常情况下应该在 100-200ms 内完成
//
// 为什么需要 10 秒（而不是 6 秒）：
//   - DNS 解析重试：如果第一次 DNS 查询失败，可能需要重试（通常 1-2 秒）
//   - 服务启动延迟：目标服务可能还在启动中，需要等待其就绪（Readiness Probe）
//   - 网络策略生效：NetworkPolicy 等配置可能需要时间生效（通常 < 5 秒）
//   - 负载均衡选择：K8s Service 选择后端 Pod 可能需要时间（通常 < 1 秒）
//
// 其他超时时间说明：
//   - 连接保持的超时：由 GRPCKeepaliveOptions() 中的 keepalive.Timeout 控制（3 秒）
//   - 请求调用的超时：由 GRPCK8sOptions() 中的 methodConfig.timeout 控制（30 秒）
//   - 如果不配置 timeout，gRPC 默认没有超时限制，会一直等待直到完成或连接断开
//   - 建议始终配置 timeout，避免请求无限期等待
//
// 建议：
//   - 如果您的环境稳定，可以降低到 6-8 秒
//   - 如果经常遇到服务启动延迟，可以保持 10-15 秒
const GRPCK8sTimeout = 10 * time.Second

// GRPCDefaultTimeout 默认的连接建立超时时间（非 K8s 环境，如本地开发）
// 本地环境通常使用 IP 地址或 localhost，DNS 解析很快，6 秒足够
const GRPCDefaultTimeout = 6 * time.Second

// TriggerGRPCReconnect 手动触发 gRPC 连接重连
// 通过检查连接状态并触发重连来强制 gRPC 重新建立连接
//
// 参数:
//   - conn: gRPC 连接对象
//
// 返回:
//   - bool: 是否需要重连（连接状态异常）
//   - error: 错误信息
//
// 性能说明:
//   - 性能消耗：极低（< 1ms）
//   - 操作：仅读取连接状态，无网络操作
//   - 适用：可以频繁调用（如每秒检查一次）
//   - 注意：如果连接状态为 Connecting，会等待最多 2 秒
//
// 使用场景:
//   - 服务恢复后，需要立即重连
//   - DNS 记录更新后，需要刷新连接
//   - 检测到连接问题时，手动触发重连
//
// 注意:
//   - 这个方法会检查连接状态，如果连接异常会返回 true
//   - 实际重连需要通过 ReconnectGRPC 或重新调用 NewClient
//
// 使用示例:
//
//	conn, _ := k.NewClient(...)
//	client := stat.NewBaseStatManageClient(conn)
//
//	// 检查连接状态，如果需要重连
//	needReconnect, _ := utils.TriggerGRPCReconnect(conn)
//	if needReconnect {
//	    // 重新建立连接
//	    newConn, err := utils.ReconnectGRPC(conn, endpoint, ...)
//	    if err != nil {
//	        log.Error("重连失败", "error", err)
//	    }
//	}
func TriggerGRPCReconnect(conn *grpc.ClientConn) (bool, error) {
	if conn == nil {
		return true, nil // 连接为空，需要重连
	}

	// 获取连接状态
	state := conn.GetState()

	// 检查连接状态，如果异常则需要重连
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return true, nil // 连接失败或已关闭，需要重连
	case connectivity.Connecting:
		// 正在连接中，等待一下
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if !conn.WaitForStateChange(ctx, state) {
			return true, nil // 等待超时，可能需要重连
		}
		return false, nil
	case connectivity.Ready:
		return false, nil // 连接正常，不需要重连
	default:
		return false, nil
	}
}

// ReconnectGRPC 使用 Kratos 重新建立 gRPC 连接
// 关闭旧连接并创建新连接，适用于需要完全重建连接的场景
//
// 参数:
//   - oldConn: 旧的 gRPC 连接对象（会被关闭）
//   - endpoint: 服务端点地址
//
// 返回:
//   - *grpc.ClientConn: 新的连接对象
//   - error: 错误信息
//
// 性能说明:
//   - 性能消耗：较高（100ms - 10秒，取决于网络状况）
//   - 操作包括：
//   - 关闭旧连接（< 1ms）
//   - DNS 解析（10-100ms，K8s 内网）
//   - TCP 连接建立（1-10ms）
//   - gRPC 握手（10-50ms）
//   - 总计：正常情况下 100-200ms，异常情况可能达到 10 秒超时
//   - 适用：仅在必要时调用（如连接失败、服务恢复）
//   - 建议：
//   - 不要频繁调用（避免连接风暴）
//   - 添加重试间隔（如指数退避）
//   - 使用连接池或单例模式管理连接
//
// 使用示例:
//
//	// 重新建立连接
//	newConn, err := utils.ReconnectGRPC(
//	    oldConn,
//	    "base-stat.micro.svc.cluster.local:10005",
//	)
//	if err != nil {
//	    log.Error("重连失败", "error", err)
//	    return
//	}
//	defer newConn.Close()
//	client := stat.NewBaseStatManageClient(newConn)
func ReconnectGRPC(oldConn *grpc.ClientConn, endpoint string) (*grpc.ClientConn, error) {
	// 关闭旧连接
	if oldConn != nil {
		_ = oldConn.Close()
	}

	// 创建新连接
	ctx, cancel := context.WithTimeout(context.Background(), GRPCK8sTimeout)
	defer cancel()

	// 使用相同的配置选项重新建立连接
	opts := MergeGRPCOptions(
		[]grpc.DialOption{grpc.WithStatsHandler(&tracing.ClientHandler{})},
		GRPCKeepaliveOptions(),
		GRPCK8sOptions(),
	)

	return k.NewClient(
		ctx,
		k.WithEndpoint(endpoint),
		k.WithMiddleware(
			recovery.Recovery(),
			metadata.Client(),
			tracing.Client(),
		),
		k.WithTimeout(GRPCK8sTimeout),
		k.WithOptions(opts...),
	)
}
