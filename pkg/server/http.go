// Package server http.go
package server

import (
	kratosjson "github.com/go-kratos/kratos/contrib/encoding/json/v3"
	"github.com/go-kratos/kratos/contrib/otel/v3/metrics"
	_ "github.com/go-kratos/kratos/v3/encoding/protojson"
	"github.com/go-kratos/kratos/v3/middleware"
	"github.com/go-kratos/kratos/v3/middleware/logging"
	"github.com/go-kratos/kratos/v3/middleware/metadata"
	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zhixingdx/oipivot-api/config"
	"github.com/zhixingdx/oipivot/backend/pkg/logger/kratoslog"
	ymetrics "github.com/zhixingdx/oipivot/backend/pkg/middleware/metrics"
	"github.com/zhixingdx/oipivot/backend/pkg/middleware/trace"
)

func init() {
	kratosjson.MarshalOptions.UseProtoNames = true
}

// CreateHTTPServer 创建采用统一中间件与默认选项的 HTTP 服务。
func CreateHTTPServer(conf *config.ServerItem, traceConf *config.Trace, logger log.Logger, ms ...middleware.Middleware) *http.Server {
	return CreateHTTPServerWithOptions(conf, traceConf, logger, nil, ms...)
}

// CreateHTTPServerWithOptions 在统一中间件基线之上注入服务级 HTTP 选项。
func CreateHTTPServerWithOptions(conf *config.ServerItem, traceConf *config.Trace, logger log.Logger, extraOptions []http.ServerOption, ms ...middleware.Middleware) *http.Server {
	if conf == nil || conf.Addr == "" {
		return nil
	}
	if conf.GetMtls() != nil {
		panic("mtls 仅用于 server.grpc")
	}
	opts := []http.ServerOption{}
	opts = append(opts, http.Address(conf.Addr))
	if conf.Network != "" {
		opts = append(opts, http.Network(conf.Network))
	}

	if conf.Timeout != nil {
		opts = append(opts, http.Timeout(conf.Timeout.AsDuration()))
	}

	use := []middleware.Middleware{}

	use = append(use, recovery.Recovery())

	if conf.IsTrace {
		use = append(use, trace.CreateMiddleware(conf, traceConf)...)
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

	use = append(use, metadata.Server())

	if len(ms) > 0 {
		use = append(use, ms...)
	}

	opts = append(opts, http.Middleware(use...))
	opts = append(opts, extraOptions...)

	srv := http.NewServer(opts...)
	if conf.IsMetrics {
		srv.Handle("/metrics", promhttp.Handler())
	}
	return srv
}
