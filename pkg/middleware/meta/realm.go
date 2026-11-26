package meta

import (
	"context"
	"regexp"
	"strings"

	"github.com/go-kratos/kratos/v3/errors"
	"github.com/go-kratos/kratos/v3/metadata"
	"github.com/zhixingdx/oipivot-api/fault"
)

const (
	// DefaultRealm 是未显式传入业务渠道时使用的兼容渠道。
	DefaultRealm = "default"
)

var realmNamePattern = regexp.MustCompile(`^[a-z][a-z0-9-]{1,30}$`)

// GetRealm 获取并规范化当前业务渠道，显式非法值必须失败而不能降级到默认渠道。
func GetRealm(ctx context.Context) (string, *errors.Error) {
	value, getErr := GetValue(ctx, realmKey)
	if getErr != nil {
		// 内部调用和未适配客户端可能没有 metadata，统一落到 default 保持兼容。
		return DefaultRealm, nil
	}
	return NormalizeRealm(value)
}

// NormalizeRealm 规范化来自请求体或 metadata 的业务渠道值。
func NormalizeRealm(value string) (string, *errors.Error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return DefaultRealm, nil
	}
	if !realmNamePattern.MatchString(value) {
		return "", fault.ErrorBadRequest("业务渠道格式非法")
	}
	return value, nil
}

// ForceInjectRealm 在派生上下文中覆盖业务渠道，并保留原有 metadata。
func ForceInjectRealm(ctx context.Context, realm string) context.Context {
	if _, ok := metadata.FromServerContext(ctx); !ok {
		ctx = metadata.NewServerContext(ctx, metadata.Metadata{})
	}
	return SetValue(ctx, realmKey, strings.TrimSpace(realm))
}
