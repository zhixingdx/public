// Package trace opentelemetry
package trace

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/zhixingdx/oipivot-api/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// SetTracerProvider _
func SetTracerProvider(conf *config.Trace) error {
	tp, err := NewTracerProvider(conf)
	if err != nil {
		return err
	}
	otel.SetTracerProvider(tp)
	return nil
}

// NewTracerProvider _
func NewTracerProvider(conf *config.Trace) (*tracesdk.TracerProvider, error) {
	serviceName := conf.Service
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(conf.Endpoint),
		otlptracegrpc.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: false,
		})),
	}
	if conf.Insecure == nil || !*conf.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	if len(conf.Headers) > 0 {
		opts = append(opts, otlptracegrpc.WithHeaders(conf.Headers))
	}

	// 认证信息通过 gRPC headers 传递，不写入 Resource（避免随 span 导出泄漏）
	auth := conf.Auth
	if auth != nil {
		if auth.Token != "" {
			opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{
				"Authorization": "Bearer " + auth.Token,
			}))
		} else if auth.Name != "" && auth.Password != "" {
			opts = append(opts, otlptracegrpc.WithHeaders(map[string]string{
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(auth.Name+":"+auth.Password)),
			}))
		}
	}

	exp, err := otlptracegrpc.New(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	rate := 1.0
	if conf.SamplingRate > 0 {
		rate = float64(conf.SamplingRate)
	}

	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(serviceName),
	}
	if conf.Env != "" {
		attrs = append(attrs, semconv.DeploymentEnvironmentKey.String(conf.Env))
	}

	return tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(rate))),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(attrs...)),
	), nil
}
