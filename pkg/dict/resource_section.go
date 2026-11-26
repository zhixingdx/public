// Package dict 资源所属分类字典
package dict

import "github.com/zhixingdx/public/pkg/lang"

// ResourceSection 资源所属分类（用于发布到校本资源的 section 字段）
type ResourceSection string

const (
	// ResourceSectionCourseware 课件
	ResourceSectionCourseware ResourceSection = "courseware"
	// ResourceSectionDesign 教学设计
	ResourceSectionDesign ResourceSection = "design"
	// ResourceSectionLesson 教案
	ResourceSectionLesson ResourceSection = "lesson"
	// ResourceSectionStudyPlan 学案
	ResourceSectionStudyPlan ResourceSection = "study_plan"
	// ResourceSectionQuestion 试题
	ResourceSectionQuestion ResourceSection = "question"
	// ResourceSectionPaper 试卷
	ResourceSectionPaper ResourceSection = "paper"
	// ResourceSectionInteractive 互动
	ResourceSectionInteractive ResourceSection = "interactive"
	// ResourceSectionImage 图片
	ResourceSectionImage ResourceSection = "image"
	// ResourceSectionMedia 音视频
	ResourceSectionMedia ResourceSection = "media"
	// ResourceSectionPack 教学包
	ResourceSectionPack ResourceSection = "pack"
	// ResourceSectionOther 其他
	ResourceSectionOther ResourceSection = "other"
)

// ResourceSectionValues 返回资源所属分类枚举值
func ResourceSectionValues() []string {
	return []string{
		string(ResourceSectionCourseware),
		string(ResourceSectionDesign),
		string(ResourceSectionLesson),
		string(ResourceSectionStudyPlan),
		string(ResourceSectionQuestion),
		string(ResourceSectionPaper),
		string(ResourceSectionInteractive),
		string(ResourceSectionImage),
		string(ResourceSectionMedia),
		string(ResourceSectionPack),
		string(ResourceSectionOther),
	}
}

// ResourceSectionLabels 返回资源所属分类中文标签（默认）
func ResourceSectionLabels() map[string]string {
	return map[string]string{
		"courseware":  "课件",
		"design":      "教学设计",
		"lesson":      "教案",
		"study_plan":  "学案",
		"question":    "试题",
		"paper":       "试卷",
		"interactive": "互动",
		"image":       "图片",
		"media":       "音视频",
		"pack":        "教学包",
		"other":       "其他",
	}
}

// ResourceSectionLangLabels 返回多语言分类标签
func ResourceSectionLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(map[string]map[string]string{
		"zh": ResourceSectionLabels(),
		"en": {
			"courseware":  "Courseware",
			"design":      "Teaching Design",
			"lesson":      "Lesson Plan",
			"study_plan":  "Study Plan",
			"question":    "Question",
			"paper":       "Paper",
			"interactive": "Interactive",
			"image":       "Image",
			"media":       "Media",
			"pack":        "Teaching Pack",
			"other":       "Other",
		},
	}, langs)
}

func (s ResourceSection) String() string {
	return string(s)
}

func (s ResourceSection) Label() string {
	return ResourceSectionLabels()[s.String()]
}

func (s ResourceSection) LangLabel(langs ...string) (string, bool) {
	labels, ok := ResourceSectionLangLabels(langs...)
	if !ok {
		return s.Label(), false
	}
	return labels[s.String()], true
}

func (s ResourceSection) Values() []string {
	return ResourceSectionValues()
}
