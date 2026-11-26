// Package dict 性别
/* cSpell:disable */
package dict

import (
	"github.com/zhixingdx/oipivot/backend/pkg/lang"
)

// Gender 性别
type Gender string

const (

	// GenderUnknown 未知
	GenderUnknown Gender = "unknown"
	// GenderMale 男
	GenderMale Gender = "male"
	// GenderFemale 女
	GenderFemale Gender = "female"
)

// Values 性别列表
func GenderValues() []string {
	return []string{
		string(GenderUnknown),
		string(GenderMale),
		string(GenderFemale),
	}
}

// GenderLabels _
func GenderLabels() map[string]string {
	return map[string]string{
		"unknown": "unknown",
		"male":    "male",
		"female":  "female",
	}
}

// 多语言支持
var genderLangMap = map[string]map[string]string{
	"zh": {
		"unknown": "未知",
		"male":    "男",
		"female":  "女",
	},
	"ru": {
		"unknown": "неизвестно",
		"male":    "мужской",
		"female":  "женский",
	},
}

func GenderLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(genderLangMap, langs)
}

// String _
func (g Gender) String() string {
	return string(g)
}

// Label _
func (g Gender) Label() string {
	return GenderLabels()[g.String()]
}

// Label _
func (g Gender) LangLabel(langs ...string) (string, bool) {
	labels, b := GenderLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}

// Values _
func (g Gender) Values() []string {
	return GenderValues()
}
