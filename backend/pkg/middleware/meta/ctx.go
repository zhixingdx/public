// Package meta ctx
package meta

import (
	"context"

	"github.com/go-kratos/kratos/v3/metadata"
	"github.com/zhixingdx/oipivot-api/fault"
)

// ValidTenantID 校验 1 到 63 字节的可打印 ASCII 标识，并拒绝首尾空格。
func ValidTenantID(id string) bool {
	if len(id) == 0 || len(id) > 63 || id[0] == ' ' || id[len(id)-1] == ' ' {
		return false
	}
	for i := 0; i < len(id); i++ {
		if id[i] < 0x20 || id[i] > 0x7e {
			return false
		}
	}
	return true
}

// GetTenantID 读取跨服务传播的唯一资源租户，拒绝重复或空的 metadata。
func GetTenantID(ctx context.Context) (string, error) {
	md, _ := metadata.FromServerContext(ctx)
	values := md.Values(TenantIDHeader)
	if len(values) != 1 || !ValidTenantID(values[0]) {
		return "", fault.ErrorBadRequest("expected one resource tenant")
	}
	return values[0], nil
}

// GetNewCtx _
func GetNewCtx(ctx context.Context) context.Context {
	newCtx := context.Background()
	md, ok := metadata.FromServerContext(ctx)
	if ok {
		// 克隆 metadata，避免新旧 ctx 共享同一个 map 引用
		newCtx = metadata.NewServerContext(newCtx, md.Clone())
	}
	return newCtx
}

// InjectTenantID 在上下文中注入租户 ID
// 用于处理 Temporal Activity 等回调场景中 ctx 可能不存在 metadata 的情况
// 如果 ctx 中已存在有效的 tenantID，则直接返回原 ctx
// 如果 ctx 中没有 metadata，则创建 metadata 并设置 tenantID
// 如果 ctx 中有 metadata 但 tenantID 为空，则在现有 metadata 上设置 tenantID
func InjectTenantID(ctx context.Context, tenantID string) context.Context {
	// 判断能否获取到租户 ID，回调 ctx 可能不存在 metadata（如 Temporal Activity）
	mTenantID, gErr := GetValue(ctx, tenantIDKey)
	if gErr == nil && mTenantID != "" {
		// 已存在有效的 tenantID，直接返回
		return ctx
	}

	// 检查 ctx 中是否存在 metadata
	if _, ok := metadata.FromServerContext(ctx); !ok {
		// ctx 中没有 metadata，需要先创建
		ctx = metadata.NewServerContext(ctx, metadata.Metadata{})
	}

	// 在现有或新建的 metadata 上设置 tenantID
	return SetValue(ctx, tenantIDKey, tenantID)
}

// ForceInjectTenantID 在上下文中强行写入租户 ID（覆盖已有值）
// 用于必须按调用方指定租户执行的内部场景（如跨租户运维、回调里重绑租户等）
// 若 ctx 无 metadata，会先创建再写入；其他 metadata 键保持不变
func ForceInjectTenantID(ctx context.Context, tenantID string) context.Context {
	if _, ok := metadata.FromServerContext(ctx); !ok {
		ctx = metadata.NewServerContext(ctx, metadata.Metadata{})
	}
	return SetValue(ctx, tenantIDKey, tenantID)
}
