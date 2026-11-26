package auth

import (
	"context"

	"github.com/go-kratos/kratos/v3/metadata"
)

// ClearHTTPIdentity 清除 HTTP 请求自报的主体身份，仅接受后续认证中间件注入的值。
func ClearHTTPIdentity(ctx context.Context) context.Context {
	if incoming, ok := metadata.FromServerContext(ctx); ok {
		// Kratos Set 忽略空字符串，必须在副本中删除，不能用空值覆盖或修改上游共享 map。
		clean := incoming.Clone()
		delete(clean, "x-md-global-staff-id")
		delete(clean, "x-md-global-user-id")
		delete(clean, "x-md-global-student-id")
		delete(clean, "x-md-global-service-id")
		ctx = metadata.NewServerContext(ctx, clean)
	}
	return ctx
}
