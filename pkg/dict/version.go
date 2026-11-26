// cSpell:disable
// Package dict 提供跨业务复用的版本状态字典。
package dict

import "github.com/zhixingdx/public/pkg/lang"

// VersionStatus 表示版本化资源唯一的生命周期状态。
type VersionStatus string

const (
	// VersionStatusDraft 表示仍允许完整替换 payload 的当前草稿。
	VersionStatusDraft VersionStatus = "draft"
	// VersionStatusPublished 表示已经发布且 payload 永久冻结的版本。
	VersionStatusPublished VersionStatus = "published"
	// VersionStatusDeprecated 表示不再允许新绑定但仍可按历史 ID 读取的版本。
	VersionStatusDeprecated VersionStatus = "deprecated"
)

// VersionStatusValues 返回版本状态的稳定枚举顺序。
func VersionStatusValues() []string {
	return []string{
		string(VersionStatusDraft),
		string(VersionStatusPublished),
		string(VersionStatusDeprecated),
	}
}

// VersionStatusLabels 返回不依赖语言环境的状态标签。
func VersionStatusLabels() map[string]string {
	return map[string]string{
		"draft":      "draft",
		"published":  "published",
		"deprecated": "deprecated",
	}
}

var versionStatusLangMap = map[string]map[string]string{
	"zh": {
		"draft":      "草稿",
		"published":  "已发布",
		"deprecated": "已废弃",
	},
	"en": {
		"draft":      "Draft",
		"published":  "Published",
		"deprecated": "Deprecated",
	},
	"ru": {
		"draft":      "Черновик",
		"published":  "Опубликовано",
		"deprecated": "Устарело",
	},
}

// VersionStatusLangLabels 返回指定语言下的版本状态标签。
func VersionStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(versionStatusLangMap, langs)
}

// String 返回状态的协议字符串。
func (v VersionStatus) String() string {
	return string(v)
}

// Label 返回状态的稳定默认标签。
func (v VersionStatus) Label() string {
	return VersionStatusLabels()[v.String()]
}

// LangLabel 返回指定语言下的状态标签，未命中时回退稳定默认标签。
func (v VersionStatus) LangLabel(langs ...string) (string, bool) {
	labels, ok := VersionStatusLangLabels(langs...)
	if !ok {
		return v.Label(), false
	}
	return labels[v.String()], true
}

// Values 返回全部允许值，供 Ent GoType 和生成器复用。
func (v VersionStatus) Values() []string {
	return VersionStatusValues()
}

// IsImmutable 判断当前版本的 payload 是否已经进入永久冻结状态。
func (v VersionStatus) IsImmutable() bool {
	return v == VersionStatusPublished || v == VersionStatusDeprecated
}
