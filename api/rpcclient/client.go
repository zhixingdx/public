// Package rpcclient 提供平台服务客户端共用的 Kratos 连接配置。
package rpcclient

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware/metadata"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	kgrpc "github.com/go-kratos/kratos/v3/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// Connect 创建带元数据传播、追踪和保活的服务连接。
func Connect(endpoint string) *grpc.ClientConn {
	conn, err := kgrpc.NewClient(context.Background(),
		kgrpc.WithEndpoint(endpoint),
		kgrpc.WithTimeout(6*time.Second),
		kgrpc.WithMiddleware(recovery.Recovery(), metadata.Client(), tracing.Client()),
		kgrpc.WithOptions(
			grpc.WithStatsHandler(&tracing.ClientHandler{}),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time: 30 * time.Second, Timeout: 3 * time.Second, PermitWithoutStream: true,
			}),
		),
	)
	if err != nil {
		panic(err)
	}
	return conn
}
