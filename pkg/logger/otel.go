// Package logger otel.go
package logger

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/zhixingdx/public/api/config"
	"github.com/zhixingdx/public/pkg/logger/kratoslog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// GetOtelLogger _
func GetOtelLogger(conf *config.Config) log.Logger {
	otelLogger, err := NewOtelLogger(conf, conf.Server)
	if err != nil {
		panic(fmt.Sprintf("init otel logger error: %v", err))
	}
	return getLogger(conf.Server, otelLogger, conf.Logger, true)
}

// OtelLogger _
type OtelLogger interface {
	log.Logger
	Close() error
}

type otelLog struct {
	loggerProvider *sdklog.LoggerProvider
	stdLogger      log.Logger
}

func (l *otelLog) Close() error {
	var errs []error
	if l.loggerProvider != nil {
		if err := l.loggerProvider.Shutdown(context.Background()); err != nil {
			errs = append(errs, fmt.Errorf("shutdown logger provider: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during shutdown: %v", errs)
	}
	return nil
}

// Log 将结构化日志发送到 OTEL，并向标准输出转发不含请求上下文的安全字段。
func (l *otelLog) Log(level log.Level, keyValues ...interface{}) error {
	// 创建日志记录
	record := otellog.Record{}
	record.SetTimestamp(time.Now())
	record.SetSeverity(convertLogLevel(level))
	record.SetBody(attribute.StringValue(extractMessage(keyValues)))

	var ctx context.Context
	stdKeyValues := make([]interface{}, 0, len(keyValues))

	// 添加属性
	kv := make(map[string]string)
	attrs := make([]attribute.KeyValue, 0, len(keyValues)/2)
	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {
			key := toString(keyValues[i])
			if key == GetCtxKey {
				oldCtx, ok := keyValues[i+1].(context.Context)
				if ok {
					ctx = oldCtx
				}
			} else {
				value := toString(keyValues[i+1])
				kv[key] = value
				attrs = append(attrs, attribute.String(key, value))
				stdKeyValues = append(stdKeyValues, keyValues[i], keyValues[i+1])
			}
		}
	}
	record.AddAttributes(attrs...)
	// 发送日志记录
	if l.loggerProvider != nil {
		logger := l.loggerProvider.Logger("logger")
		if ctx != nil {
			logger.Emit(ctx, record)
		} else {
			logger.Emit(context.Background(), record)
		}
		if l.stdLogger != nil {
			// get_ctx 只用于关联 OTEL trace；转发到标准输出会序列化完整请求上下文并泄露认证凭证。
			_ = l.stdLogger.Log(level, stdKeyValues...)
		}
	}

	return nil
}

// convertLogLevel 将 go-kratos 日志级别转换为 OpenTelemetry 日志级别
func convertLogLevel(level log.Level) otellog.Severity {
	switch level {
	case log.LevelDebug:
		return otellog.SeverityDebug
	case log.LevelInfo:
		return otellog.SeverityInfo
	case log.LevelWarn:
		return otellog.SeverityWarn
	case log.LevelError, log.LevelFatal:
		return otellog.SeverityError
	default:
		return otellog.SeverityInfo
	}
}

// extractMessage 从键值对中提取消息
func extractMessage(keyValues []interface{}) string {
	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {
			key := toString(keyValues[i])
			if key == "msg" || key == "message" || key == "operation" {
				return toString(keyValues[i+1])
			}
		}
	}
	return "日志记录"
}

// NewOtelLogger _
func NewOtelLogger(conf *config.Config, serverConfig *config.Server) (OtelLogger, error) {
	var loggerProvider *sdklog.LoggerProvider
	var res *resource.Resource
	var err error

	// 创建资源，设置服务信息（deployment.environment 与 trace/metrics 一致）
	resAttrs := []attribute.KeyValue{
		semconv.ServiceName(serverConfig.Name),
		semconv.ServiceVersion(serverConfig.Version),
		semconv.ServiceInstanceID(serverConfig.Id),
		semconv.ContainerID(serverConfig.Id),
		semconv.HostName(serverConfig.Id),
	}
	if serverConfig.Env != "" {
		resAttrs = append(resAttrs, semconv.DeploymentEnvironmentKey.String(serverConfig.Env))
	}
	res, err = resource.New(context.Background(), resource.WithAttributes(resAttrs...))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	opts := conf.Logger
	// 获取 endpoint
	endpoint := opts.Endpoint
	if endpoint == "" {
		return nil, fmt.Errorf("Logger.Endpoint must be specified")
	}

	// 准备 OTLP 日志导出器选项
	otlpOptions := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(endpoint),
		otlploggrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
		otlploggrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: false,
		})),
	}
	if opts.Insecure == nil || !*opts.Insecure {
		otlpOptions = append(otlpOptions, otlploggrpc.WithInsecure())
	}
	if len(opts.Headers) > 0 {
		otlpOptions = append(otlpOptions, otlploggrpc.WithHeaders(opts.Headers))
	}

	// 添加认证配置
	if opts.Token != "" {
		// 使用 token 认证
		headers := map[string]string{
			"Authorization": "Bearer " + opts.Token,
		}
		otlpOptions = append(otlpOptions, otlploggrpc.WithHeaders(headers))
	} else if opts.Username != "" && opts.Password != "" {
		// 使用用户名密码认证
		headers := map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(opts.Username+":"+opts.Password)),
		}
		otlpOptions = append(otlpOptions, otlploggrpc.WithHeaders(headers))
	}

	// 创建日志导出器
	logExporter, err := otlploggrpc.New(
		context.Background(),
		otlpOptions...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP log exporter: %w", err)
	}

	// 创建日志处理器
	processor := sdklog.NewBatchProcessor(logExporter)

	// 创建日志提供者
	loggerProvider = sdklog.NewLoggerProvider(
		sdklog.WithProcessor(processor),
		sdklog.WithResource(res),
	)

	otelLog := &otelLog{
		loggerProvider: loggerProvider,
	}
	if opts.AlsoStd {
		stdLogger := GetStdLogger(conf)
		otelLog.stdLogger = stdLogger
	}
	return otelLog, nil
}

// // BuildCtxWithOTEL 构建一个包含链路追踪 ID 的 Context
// // 符合 OpenTelemetry 的链路追踪规范
// func BuildCtxWithOTEL(kv map[string]string) context.Context {
// 	ctx := context.Background()
// 	spanContext := trace.NewSpanContext(trace.SpanContextConfig{})
// 	// 使用 OpenTelemetry 的标准属性键
// 	if traceID, exists := kv["trace_id"]; exists && traceID != "" {
// 		if tid, err := trace.TraceIDFromHex(traceID); err == nil {
// 			spanContext = spanContext.WithTraceID(tid)
// 		}
// 		if spanID, exists := kv["span_id"]; exists && spanID != "" {
// 			if sid, err := trace.SpanIDFromHex(spanID); err == nil {
// 				spanContext = spanContext.WithSpanID(sid)
// 			}
// 		}
// 	}
// 	ctx = trace.ContextWithSpanContext(ctx, spanContext)
// 	return ctx
// }
