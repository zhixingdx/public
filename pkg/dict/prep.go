// Package dict 备课域字典
package dict

import "github.com/zhixingdx/public/pkg/lang"

// PrepBookType 备课本类型
type PrepBookType string

const (
	// PrepBookTypeSyncLesson 同步课（关联教材章节）
	PrepBookTypeSyncLesson PrepBookType = "syncLesson"
	// PrepBookTypeSpecialTopic 专题课（不强制教材）
	PrepBookTypeSpecialTopic PrepBookType = "specialTopic"
)

// PrepBookTypeValues 返回备课本类型枚举值
func PrepBookTypeValues() []string {
	return []string{
		string(PrepBookTypeSyncLesson),
		string(PrepBookTypeSpecialTopic),
	}
}

// PrepBookTypeLabels 返回备课本类型中文标签
func PrepBookTypeLabels() map[string]string {
	return map[string]string{
		"syncLesson":   "同步课",
		"specialTopic": "专题课",
	}
}

// PrepBookTypeLangLabels 返回多语言标签
func PrepBookTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(map[string]map[string]string{
		"zh": PrepBookTypeLabels(),
		"en": {
			"syncLesson":   "Sync Lesson",
			"specialTopic": "Special Topic",
		},
	}, langs)
}

func (t PrepBookType) String() string {
	return string(t)
}

func (t PrepBookType) Label() string {
	return PrepBookTypeLabels()[t.String()]
}

func (t PrepBookType) LangLabel(langs ...string) (string, bool) {
	labels, ok := PrepBookTypeLangLabels(langs...)
	if !ok {
		return t.Label(), false
	}
	return labels[t.String()], true
}

func (t PrepBookType) Values() []string {
	return PrepBookTypeValues()
}

// PrepPackageDefaultLanguage 备课包默认语言
type PrepPackageDefaultLanguage string

const (
	PrepPackageDefaultLanguageZh   PrepPackageDefaultLanguage = "zh"
	PrepPackageDefaultLanguageEn   PrepPackageDefaultLanguage = "en"
	PrepPackageDefaultLanguageZhEn PrepPackageDefaultLanguage = "zh-en"
)

// PrepPackageDefaultLanguageValues 返回默认语言枚举值
func PrepPackageDefaultLanguageValues() []string {
	return []string{
		string(PrepPackageDefaultLanguageZh),
		string(PrepPackageDefaultLanguageEn),
		string(PrepPackageDefaultLanguageZhEn),
	}
}

// PrepPackageDefaultLanguageLabels 返回默认语言中文标签
func PrepPackageDefaultLanguageLabels() map[string]string {
	return map[string]string{
		"zh":    "中文",
		"en":    "英文",
		"zh-en": "中英双语",
	}
}

func (l PrepPackageDefaultLanguage) String() string {
	return string(l)
}

func (l PrepPackageDefaultLanguage) Label() string {
	return PrepPackageDefaultLanguageLabels()[l.String()]
}

func (l PrepPackageDefaultLanguage) Values() []string {
	return PrepPackageDefaultLanguageValues()
}

// PrepPackageFileKind 备课包内文件类型
type PrepPackageFileKind string

const (
	PrepPackageFileKindTeachingDesign   PrepPackageFileKind = "teachingDesign"
	PrepPackageFileKindCourseware       PrepPackageFileKind = "courseware"
	PrepPackageFileKindTranscript       PrepPackageFileKind = "transcript"
	PrepPackageFileKindStudentTaskSheet PrepPackageFileKind = "studentTaskSheet"
	PrepPackageFileKindTieredHomework   PrepPackageFileKind = "tieredHomework"
	PrepPackageFileKindRubric           PrepPackageFileKind = "rubric"
	PrepPackageFileKindLessonPlan       PrepPackageFileKind = "lessonPlan"
	PrepPackageFileKindLectureScript    PrepPackageFileKind = "lectureScript"
	PrepPackageFileKindTeachingReflection PrepPackageFileKind = "teachingReflection"
	PrepPackageFileKindOther            PrepPackageFileKind = "other"
)

// PrepPackageFileKindValues 返回备课包文件类型枚举值
func PrepPackageFileKindValues() []string {
	return []string{
		string(PrepPackageFileKindTeachingDesign),
		string(PrepPackageFileKindCourseware),
		string(PrepPackageFileKindTranscript),
		string(PrepPackageFileKindStudentTaskSheet),
		string(PrepPackageFileKindTieredHomework),
		string(PrepPackageFileKindRubric),
		string(PrepPackageFileKindLessonPlan),
		string(PrepPackageFileKindLectureScript),
		string(PrepPackageFileKindTeachingReflection),
		string(PrepPackageFileKindOther),
	}
}

// PrepPackageFileKindLabels 返回备课包文件类型中文标签
func PrepPackageFileKindLabels() map[string]string {
	return map[string]string{
		"teachingDesign":     "教学设计",
		"courseware":           "课件",
		"transcript":           "逐字稿",
		"studentTaskSheet":     "学生任务单",
		"tieredHomework":       "分层作业",
		"rubric":               "评价量规",
		"lessonPlan":           "教案",
		"lectureScript":        "说课稿",
		"teachingReflection":   "教学反思",
		"other":                "其他",
	}
}

func (k PrepPackageFileKind) String() string {
	return string(k)
}

func (k PrepPackageFileKind) Label() string {
	return PrepPackageFileKindLabels()[k.String()]
}

func (k PrepPackageFileKind) Values() []string {
	return PrepPackageFileKindValues()
}
