// Package dict 校本资源来源类型字典
package dict

import "github.com/zhixingdx/public/pkg/lang"

// SchoolResourceSource 校本资源来源类型
type SchoolResourceSource string

const (
	// SchoolResourceSourcePersonal 个人资源发布
	SchoolResourceSourcePersonal SchoolResourceSource = "personal"
	// SchoolResourceSourceAdmin 管理员上传
	SchoolResourceSourceSchool SchoolResourceSource = "school"
	// SchoolResourceSourcePlatform 平台资源
	SchoolResourceSourcePlatform SchoolResourceSource = "platform"
)

// SchoolResourceSourceValues 返回校本资源来源枚举值
func SchoolResourceSourceValues() []string {
	return []string{
		string(SchoolResourceSourcePersonal),
		string(SchoolResourceSourceSchool),
		string(SchoolResourceSourcePlatform),
	}
}

// SchoolResourceSourceLabels 返回校本资源来源中文标签
func SchoolResourceSourceLabels() map[string]string {
	return map[string]string{
		"personal": "个人发布",
		"school":   "学校上传",
		"platform": "平台资源",
	}
}

// SchoolResourceSourceLangLabels 返回多语言标签
func SchoolResourceSourceLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(map[string]map[string]string{
		"zh": SchoolResourceSourceLabels(),
		"en": {
			"personal": "Personal Publish",
			"school":   "School Upload",
			"platform": "Platform",
		},
	}, langs)
}

func (s SchoolResourceSource) String() string {
	return string(s)
}

func (s SchoolResourceSource) Label() string {
	return SchoolResourceSourceLabels()[s.String()]
}

func (s SchoolResourceSource) LangLabel(langs ...string) (string, bool) {
	labels, ok := SchoolResourceSourceLangLabels(langs...)
	if !ok {
		return s.Label(), false
	}
	return labels[s.String()], true
}

func (s SchoolResourceSource) Values() []string {
	return SchoolResourceSourceValues()
}
