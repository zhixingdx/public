// Package dict 文件相关字典
package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

// FileCategory 文件大类（用于资源库筛选）
type FileCategory string

const (
	// FileCategoryDocument 文档（PDF、Office、纯文本等）
	FileCategoryDocument FileCategory = "document"
	// FileCategoryMedia 音视频
	FileCategoryMedia FileCategory = "media"
	// FileCategoryImage 图片
	FileCategoryImage FileCategory = "image"
	// FileCategoryOther 其他（含文件夹、未知 MIME）
	FileCategoryOther FileCategory = "other"
)

// FileCategoryValues 返回所有文件大类的枚举值
func FileCategoryValues() []string {
	return []string{
		string(FileCategoryDocument),
		string(FileCategoryMedia),
		string(FileCategoryImage),
		string(FileCategoryOther),
	}
}

// FileCategoryLabels 返回文件大类中文标签（默认）
func FileCategoryLabels() map[string]string {
	return map[string]string{
		"document": "文档",
		"media":    "音视频",
		"image":    "图片",
		"other":    "其他",
	}
}

// FileCategoryLangLabels 多语言标签
func FileCategoryLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(map[string]map[string]string{
		"zh": FileCategoryLabels(),
		"en": {
			"document": "Document",
			"media":    "Media",
			"image":    "Image",
			"other":    "Other",
		},
	}, langs)
}

func (c FileCategory) String() string {
	return string(c)
}

func (c FileCategory) Label() string {
	return FileCategoryLabels()[c.String()]
}

func (c FileCategory) LangLabel(langs ...string) (string, bool) {
	labels, ok := FileCategoryLangLabels(langs...)
	if !ok {
		return c.Label(), false
	}
	return labels[c.String()], true
}

func (c FileCategory) Values() []string {
	return FileCategoryValues()
}
