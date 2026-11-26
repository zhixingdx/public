// Package dict 租户相关字典
/* cSpell:disable */
package dict

import (
	"github.com/zhixingdx/oipivot/backend/pkg/lang"
)

// TenantVisibilityScope 租户可见范围
// 用于控制资源在不同租户之间的可见性
type TenantVisibilityScope string

const (
	// TenantVisibilityScopeSelf 个人可见（仅所属租户内的本人可见，Schema 须显式启用）。
	TenantVisibilityScopeSelf TenantVisibilityScope = "self"
	// TenantVisibilityScopeInternal 内部可见（仅本租户可见）
	TenantVisibilityScopeInternal TenantVisibilityScope = "internal"
	// TenantVisibilityScopeCross 跨租户可见（指定租户可见）
	TenantVisibilityScopeCross TenantVisibilityScope = "cross"
	// TenantVisibilityScopePublic 公开可见（所有租户可见）
	TenantVisibilityScopePublic TenantVisibilityScope = "public"
)

// TenantVisibilityScopeValues 返回默认范围，个人范围由具备本人鉴权的 Schema 显式追加。
func TenantVisibilityScopeValues() []string {
	return []string{
		string(TenantVisibilityScopeInternal),
		string(TenantVisibilityScopeCross),
		string(TenantVisibilityScopePublic),
	}
}

// TenantVisibilityScopeLabels 租户可见范围标签
func TenantVisibilityScopeLabels() map[string]string {
	return map[string]string{
		"self":     "self",
		"internal": "internal",
		"cross":    "cross",
		"public":   "public",
	}
}

// 多语言支持
var tenantVisibilityScopeLangMap = map[string]map[string]string{
	"zh": {
		"self":     "仅本人可见",
		"internal": "内部可见",
		"cross":    "跨单位可见",
		"public":   "公开可见",
	},
	"en": {
		"self":     "Self Only",
		"internal": "Internal Only",
		"cross":    "Cross Organization",
		"public":   "Public",
	},
	"ru": {
		"self":     "только для себя",
		"internal": "внутренний",
		"cross":    "межорганизационный",
		"public":   "публичный",
	},
}

// TenantVisibilityScopeLangLabels 获取多语言标签
func TenantVisibilityScopeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(tenantVisibilityScopeLangMap, langs)
}

// String 转换为字符串
func (t TenantVisibilityScope) String() string {
	return string(t)
}

// Label 获取标签
func (t TenantVisibilityScope) Label() string {
	return TenantVisibilityScopeLabels()[t.String()]
}

// LangLabel 获取多语言标签
func (t TenantVisibilityScope) LangLabel(langs ...string) (string, bool) {
	labels, b := TenantVisibilityScopeLangLabels(langs...)
	if !b {
		return t.Label(), false
	}
	return labels[t.String()], true
}

// Values 保持既有 Schema 的默认枚举集合，不自动开启个人范围。
func (t TenantVisibilityScope) Values() []string {
	return TenantVisibilityScopeValues()
}

// IsValid 仅校验默认三种范围，个人范围由显式启用它的 Schema 校验。
func (t TenantVisibilityScope) IsValid() bool {
	switch t {
	case TenantVisibilityScopeInternal, TenantVisibilityScopeCross, TenantVisibilityScopePublic:
		return true
	}
	return false
}
