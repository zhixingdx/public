// Package server grpc.go
package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	"github.com/go-kratos/kratos/contrib/otel/v3/tracing"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/logging"
	"github.com/go-kratos/kratos/v3/middleware/metadata"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/grpc"
	"github.com/zhixingdx/oipivot-api/config"
	"github.com/zhixingdx/oipivot/backend/pkg/logger/kratoslog"
	ymetrics "github.com/zhixingdx/oipivot/backend/pkg/middleware/metrics"
	"github.com/zhixingdx/oipivot/backend/pkg/utils"
	"go.opentelemetry.io/otel"
)

// CreateGRPCServer 使用仓库统一中间件创建默认 gRPC Server。
func CreateGRPCServer(conf *config.ServerItem, logger log.Logger, ms ...middleware.Middleware) *grpc.Server {
	return CreateGRPCServerWithOptions(conf, logger, nil, ms...)
}

// CreateGRPCServerWithOptions 创建支持显式 Kratos ServerOption 的 gRPC Server。
func CreateGRPCServerWithOptions(
	conf *config.ServerItem,
	logger log.Logger,
	serverOptions []grpc.ServerOption,
	ms ...middleware.Middleware,
) *grpc.Server {
	if conf == nil || conf.Addr == "" {
		return nil
	}

	opts := append([]grpc.ServerOption(nil), serverOptions...)
	if files := conf.GetMtls(); files != nil {
		tlsConfig, err := loadGRPCMutualTLS(files)
		if err != nil {
			panic(fmt.Errorf("server.grpc.mtls: %w", err))
		}
		opts = append(opts, grpc.TLSConfig(tlsConfig))
	}
	if conf.Network != "" {
		opts = append(opts, grpc.Network(conf.Network))
	}

	if conf.Addr != "" {
		opts = append(opts, grpc.Address(conf.Addr))
	}

	if conf.Timeout != nil {
		opts = append(opts, grpc.Timeout(conf.Timeout.AsDuration()))
	}

	use := []middleware.Middleware{}

	use = append(use, recovery.Recovery())

	if conf.IsTrace {
		use = append(use, tracing.Server(tracing.WithTracerProvider(otel.GetTracerProvider())))
	}

	if conf.IsLog {
		use = append(use, logging.Server(log.AsSlog(logger)))
	}

	if conf.IsMetrics {
		metricRequests, metricSeconds := ymetrics.GetDefaultMetrics(nil)
		use = append(use, metrics.Server(
			metrics.WithSeconds(metricSeconds),
			metrics.WithRequests(metricRequests),
		))
	}

	//  validate.Validator()
	use = append(use, metadata.Server())
	if len(ms) > 0 {
		use = append(use, ms...)
	}

	opts = append(opts, grpc.Middleware(use...))

	// 添加服务端 keepalive 配置，允许客户端发送 ping，避免 "too_many_pings" 错误
	// 使用 grpc.Options 传递底层的 google.golang.org/grpc.ServerOption
	opts = append(opts, grpc.Options(utils.GRPCServerKeepaliveOptions()...))

	srv := grpc.NewServer(opts...)

	return srv
}

func loadGRPCMutualTLS(files *config.MutualTLS) (*tls.Config, error) {
	if files.CertFile == "" || files.KeyFile == "" || files.ClientCAFile == "" {
		return nil, fmt.Errorf("需要 certFile、keyFile 和 clientCAFile")
	}
	cert, err := tls.LoadX509KeyPair(files.CertFile, files.KeyFile)
	if err != nil {
		return nil, err
	}
	pem, err := os.ReadFile(files.ClientCAFile)
	if err != nil {
		return nil, err
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(pem) {
		return nil, fmt.Errorf("无效的客户端 CA")
	}
	return &tls.Config{MinVersion: tls.VersionTLS13, Certificates: []tls.Certificate{cert}, ClientCAs: roots, ClientAuth: tls.RequireAndVerifyClientCert}, nil
}
