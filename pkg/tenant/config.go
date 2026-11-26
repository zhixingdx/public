// Package tenant 多租户配置工具包
// 提供 SaaS 系统中多租户配置的通用模式实现
package tenant

import "github.com/samber/lo"

// Config 多租户配置通用结构（泛型）
// T 可以是任意配置类型，如 []*WorkflowDefinitionConfig, *SMSConfig 等
//
// 配置优先级：
//  1. TenantConfigs[tenantID] - 租户自定义配置（优先级最高）
//  2. 租户列表控制 - PlatformTenants / NonPlatformTenants
//  3. 符合名单策略时选择 Platform，否则返回 SourceNone
//
// 使用示例：
//
//	type CustomConfig struct {
//	    Definitions *tenant.Config[[]*WorkflowDefinitionConfig] `json:"definitions"`
//	    SMS         *tenant.Config[*SMSConfig]                  `json:"sms"`
//	}
type Config[T any] struct {
	// 平台配置（按 UsePlatform 和名单选择）
	Platform T `json:"platform"`
	// 租户自定义配置（优先级最高）
	TenantConfigs map[string]T `json:"tenant_configs"`
	// 默认是否使用平台配置
	UsePlatform bool `json:"use_platform"`
	// 使用平台配置的租户列表（UsePlatform=false 时生效）
	PlatformTenants []string `json:"platform_tenants"`
	// 不使用平台配置的租户列表（UsePlatform=true 时生效）
	NonPlatformTenants []string `json:"non_platform_tenants"`
}

// ConfigResult 配置获取结果
type ConfigResult[T any] struct {
	// 配置值
	Value T
	// 是否找到配置
	Found bool
	// 配置来源：platform=平台配置, tenant=租户自定义, none=未找到
	Source ConfigSource
}

// ConfigSource 配置来源
type ConfigSource string

const (
	// SourcePlatform 平台配置
	SourcePlatform ConfigSource = "platform"
	// SourceTenant 租户自定义
	SourceTenant ConfigSource = "tenant"
	// SourceNone 未找到配置（需从数据库获取）
	SourceNone ConfigSource = "none"
)

// Get 获取租户配置
// 返回配置结果，包含配置值、是否找到、配置来源
func (c *Config[T]) Get(tenantID string) ConfigResult[T] {
	if c == nil {
		var zero T
		return ConfigResult[T]{Value: zero, Found: false, Source: SourceNone}
	}

	// 1. 优先查租户自定义配置
	if c.TenantConfigs != nil {
		if value, ok := c.TenantConfigs[tenantID]; ok {
			return ConfigResult[T]{Value: value, Found: true, Source: SourceTenant}
		}
	}

	// 2. 判断是否使用平台配置
	if c.ShouldUsePlatform(tenantID) {
		return ConfigResult[T]{Value: c.Platform, Found: true, Source: SourcePlatform}
	}

	// 3. 不使用平台配置且无租户自定义，返回未找到
	var zero T
	return ConfigResult[T]{Value: zero, Found: false, Source: SourceNone}
}

// GetValue 简化版获取，只返回配置值和是否找到
func (c *Config[T]) GetValue(tenantID string) (T, bool) {
	result := c.Get(tenantID)
	return result.Value, result.Found
}

// ShouldUsePlatform 判断租户是否应使用平台配置
func (c *Config[T]) ShouldUsePlatform(tenantID string) bool {
	if c == nil {
		return false
	}

	// 有租户自定义配置时，不使用平台配置
	if c.TenantConfigs != nil {
		if _, ok := c.TenantConfigs[tenantID]; ok {
			return false
		}
	}

	if c.UsePlatform {
		// 默认使用平台配置，除非在 NonPlatformTenants 列表中
		return !lo.Contains(c.NonPlatformTenants, tenantID)
	}
	// 默认不使用平台配置，除非在 PlatformTenants 列表中
	return lo.Contains(c.PlatformTenants, tenantID)
}

// HasTenantConfig 检查租户是否有自定义配置
func (c *Config[T]) HasTenantConfig(tenantID string) bool {
	if c == nil || c.TenantConfigs == nil {
		return false
	}
	_, ok := c.TenantConfigs[tenantID]
	return ok
}

// ListTenants 列出所有有自定义配置的租户
func (c *Config[T]) ListTenants() []string {
	if c == nil || c.TenantConfigs == nil {
		return nil
	}
	tenants := make([]string, 0, len(c.TenantConfigs))
	for tenantID := range c.TenantConfigs {
		tenants = append(tenants, tenantID)
	}
	return tenants
}
