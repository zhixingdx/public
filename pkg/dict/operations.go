// Package dict 运营相关字典定义
/* cSpell:disable */
package dict

import (
	"github.com/zhixingdx/oipivot/backend/pkg/lang"
)

// AnnouncementStatus 公告状态
type AnnouncementStatus string

const (
	// 公告状态：未上架
	AnnouncementStatusUnshelved AnnouncementStatus = "unshelved"
	// 公告状态：已上架
	AnnouncementStatusShelved AnnouncementStatus = "shelved"
	// 公告状态：已下架
	AnnouncementStatusOffShelf AnnouncementStatus = "offShelf"
)

// AnnouncementStatusValues 公告状态列表
func AnnouncementStatusValues() []string {
	return []string{
		string(AnnouncementStatusUnshelved),
		string(AnnouncementStatusShelved),
		string(AnnouncementStatusOffShelf),
	}
}

// AnnouncementStatusLabels 公告状态标签映射
func AnnouncementStatusLabels() map[string]string {
	return map[string]string{
		"unshelved": "unshelved",
		"shelved":   "shelved",
		"offShelf":  "offShelf",
	}
}

// 多语言支持
var announcementStatusLangMap = map[string]map[string]string{
	"zh": {
		"unshelved": "未上架",
		"shelved":   "已上架",
		"offShelf":  "已下架",
	},
	"ru": {
		"unshelved": "не выставлен",
		"shelved":   "выставлен",
		"offShelf":  "снят с полки",
	},
}

// AnnouncementStatusLangLabels 获取多语言标签
func AnnouncementStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(announcementStatusLangMap, langs)
}

// String 返回字符串值
func (a AnnouncementStatus) String() string {
	return string(a)
}

// Label 返回标签
func (a AnnouncementStatus) Label() string {
	return AnnouncementStatusLabels()[a.String()]
}

// LangLabel 返回多语言标签
func (a AnnouncementStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := AnnouncementStatusLangLabels(langs...)
	if !b {
		return a.Label(), false
	}
	return labels[a.String()], true
}

// Values 返回所有状态值列表
func (a AnnouncementStatus) Values() []string {
	return AnnouncementStatusValues()
}
