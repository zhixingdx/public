// Package logger logger.go
package logger

import (
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/zhixingdx/public/api/config"
	"github.com/zhixingdx/public/pkg/logger/kratoslog"
)

var (
	// defaultSensitiveKeys 覆盖所有服务中必须无条件脱敏的凭据字段。
	defaultSensitiveKeys = []string{"clientSecret", "client_secret"}
	// defaultSensitiveRegex 覆盖 Kratos args 中的 protobuf text 与 JSON 两种嵌套表示。
	defaultSensitiveRegex = []string{
		`clientSecret:"(\\.|[^"\\])*"`,
		`client_secret:"(\\.|[^"\\])*"`,
		`"clientSecret":"(\\.|[^"\\])*"`,
		`"client_secret":"(\\.|[^"\\])*"`,
	}
)

const (
	GetCtxKey = "get_ctx"
	// grpcHealthCheckOperation 与 gRPC 官方 Health 服务的 Check 方法全名一致，用于过滤探针访问日志
	grpcHealthCheckOperation = "/grpc.health.v1.Health/Check"
	// updateKubernetesSecretOperation 的请求包含任意 Secret 正文，不能依赖字段名或正则逐项脱敏。
	updateKubernetesSecretOperation = "/api.monolith.devops.MonolithDevopsManagePlus/UpdateKubernetesSecret"
)

// redactSensitiveOperationArgs 保留操作与状态日志，但隐藏直接承载任意私密正文的请求参数。
func redactSensitiveOperationArgs(_ log.Level, keyvals ...any) bool {
	secretOperation := false
	for i := 0; i+1 < len(keyvals); i += 2 {
		if keyvals[i] == "operation" && keyvals[i+1] == updateKubernetesSecretOperation {
			secretOperation = true
			break
		}
	}
	if secretOperation {
		for i := 0; i+1 < len(keyvals); i += 2 {
			if keyvals[i] == "args" {
				keyvals[i+1] = "***"
			}
		}
	}
	return false
}

// shouldSkipGRPCHealthCheckAccessLog 匹配 logging.Server 中间件打印的访问日志：跳过健康检查，避免探针刷屏
func shouldSkipGRPCHealthCheckAccessLog(_ log.Level, keyvals ...any) bool {
	for i := 0; i+1 < len(keyvals); i += 2 {
		k, ok := keyvals[i].(string)
		if !ok || k != "operation" {
			continue
		}
		v, ok := keyvals[i+1].(string)
		return ok && v == grpcHealthCheckOperation
	}
	return false
}

func getCtx() log.Valuer {
	return func(ctx context.Context) any {
		return ctx
	}
}

// compileSensitiveRegex 忽略无效的可选配置，并在 logger 初始化时一次性编译有效规则。
func compileSensitiveRegex(patterns []string) []*regexp.Regexp {
	regexps := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			regexps = append(regexps, re)
		}
	}
	return regexps
}

// redactSensitiveLogValues 在日志值内部替换敏感片段，并保留字段名供排障识别。
func redactSensitiveLogValues(regexps []*regexp.Regexp) log.FilterOption {
	return log.FilterFunc(func(_ log.Level, keyvals ...any) bool {
		for i := 1; i < len(keyvals); i += 2 {
			original := toString(keyvals[i])
			value := original
			for _, re := range regexps {
				value = re.ReplaceAllStringFunc(value, func(match string) string {
					if colonIndex := strings.Index(match, ":"); colonIndex > 0 {
						return match[:colonIndex] + ":\"***\""
					}
					return "***"
				})
			}
			if value != original {
				keyvals[i] = value
			}
		}
		return false
	})
}

// GetLogger _
func GetLogger(conf *config.Config) log.Logger {
	logger := conf.Logger
	if logger == nil {
		return GetStdLogger(conf)
	}
	if logger.Provider == "tencent" {
		return GetTencentLogger(conf)
	}
	if logger.Provider == "otel" {
		return GetOtelLogger(conf)
	}
	return GetStdLogger(conf)
}

// GetStdLogger _
func GetStdLogger(conf *config.Config) log.Logger {
	logger := log.NewStdLogger(os.Stdout)
	return getLogger(conf.Server, logger, conf.Logger, false)
}

// getLogger 组合服务日志上下文、配置过滤规则与全局敏感字段保护。
func getLogger(service *config.Server, logger log.Logger, loggerConf *config.Logger, isWithCtx bool) log.Logger {
	var filterOptions []log.FilterOption
	filterKeys := append([]string(nil), defaultSensitiveKeys...)
	sensitivePatterns := append([]string(nil), defaultSensitiveRegex...)

	if loggerConf != nil {
		// 按日志级别过滤
		if loggerConf.FilterLevel != "" {
			level := log.ParseLevel(loggerConf.FilterLevel)
			filterOptions = append(filterOptions, log.FilterLevel(level))
		}

		// 按key过滤（脱敏）
		filterKeys = append(filterKeys, loggerConf.FilterKeys...)

		// 按value过滤（脱敏）
		if len(loggerConf.FilterValues) > 0 {
			filterOptions = append(filterOptions, log.FilterValue(loggerConf.FilterValues...))
		}

		// 按敏感信息正则表达式脱敏
		sensitivePatterns = append(sensitivePatterns, loggerConf.SensitiveRegex...)
	}

	filterOptions = append(filterOptions, log.FilterKey(filterKeys...))
	filterOptions = append(filterOptions, log.FilterFunc(redactSensitiveOperationArgs))
	if regexps := compileSensitiveRegex(sensitivePatterns); len(regexps) > 0 {
		filterOptions = append(filterOptions, redactSensitiveLogValues(regexps))
	}

	// 默认不输出 gRPC Health/Check 的访问日志（logging.Server 中间件中的 operation）
	filterOptions = append(filterOptions, log.FilterFunc(shouldSkipGRPCHealthCheckAccessLog))

	if len(filterOptions) > 0 {
		logger = log.NewFilter(logger, filterOptions...)
	}

	kvs := []any{
		"ts", log.DefaultTimestamp,
		"caller", log.Caller(8),
		"service.id", service.Id,
		"service.name", service.Name,
		"service.version", service.Version,
		"deployment.environment", service.Env,
	}
	if isWithCtx {
		kvs = append(kvs, GetCtxKey, getCtx())
	}

	return log.With(logger, kvs...)
}
