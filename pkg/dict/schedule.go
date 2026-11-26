// Package dict 性别
/* cSpell:disable */
package dict

import (
	"github.com/zhixingdx/public/pkg/lang"
)

// PublicScope 公开范围
type PublicScope string

const (
	// 互联网公开
	PublicScopeOnline PublicScope = "public"
	// 本单位
	PublicScopeUnit PublicScope = "unit"
	// 本单位及其他单位
	PublicScopeUnitAndOther PublicScope = "unit_and_other"
)

// Values 公开范围列表
func PublicScopeValues() []string {
	return []string{
		string(PublicScopeOnline),
		string(PublicScopeUnit),
		string(PublicScopeUnitAndOther),
	}
}

// PublicScopeLabels _
func PublicScopeLabels() map[string]string {
	return map[string]string{
		"public":         "public",
		"unit":           "unit",
		"unit_and_other": "unit_and_other",
	}
}

// 多语言支持
var publicScopeLangMap = map[string]map[string]string{
	"zh": {
		"public":         "互联网公开",
		"unit":           "本单位",
		"unit_and_other": "本单位及其他单位",
	},
	"ru": {
		"public":         "онлайн",
		"unit":           "юнит",
		"unit_and_other": "юнит и другие",
	},
}

func PublicScopeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(publicScopeLangMap, langs)
}

// String _
func (g PublicScope) String() string {
	return string(g)
}

// Label _
func (g PublicScope) Label() string {
	return PublicScopeLabels()[g.String()]
}

// Label _
func (g PublicScope) LangLabel(langs ...string) (string, bool) {
	labels, b := PublicScopeLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}

// Values _
func (g PublicScope) Values() []string {
	return PublicScopeValues()
}

type LessonType string

const (
	// 课堂类型：直播
	LessonTypeLive LessonType = "live"
	// 课堂类型：互动
	LessonTypeInteractive LessonType = "interactive"
	// 课堂类型：互动直播
	LessonTypeLiveInteractive LessonType = "live_interactive"
	// 其他：其他
	LessonTypeOther LessonType = "other"
)

// Values 课程类型列表
func LessonTypeValues() []string {
	return []string{
		string(LessonTypeLive),
		string(LessonTypeInteractive),
		string(LessonTypeLiveInteractive),
		string(LessonTypeOther),
	}
}

// LessonTypeLabels _
func LessonTypeLabels() map[string]string {
	return map[string]string{
		"live":             "live",
		"interactive":      "interactive",
		"live_interactive": "live_interactive",
		"other":            "other",
	}
}

// 多语言支持
var LessonTypeLangMap = map[string]map[string]string{
	"zh": {
		"live":             "直播课堂",
		"interactive":      "直播课堂",
		"live_interactive": "互动直播",
		"other":            "其他",
	},
	"ru": {
		"live":             "жить",
		"interactive":      "интерактивный",
		"live_interactive": "живой интерактивный",
		"other":            "другой",
	},
}

func LessonTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(LessonTypeLangMap, langs)
}

// String _
func (g LessonType) String() string {
	return string(g)
}

// Label _
func (g LessonType) Label() string {
	return LessonTypeLabels()[g.String()]
}

// Label _
func (g LessonType) LangLabel(langs ...string) (string, bool) {
	labels, b := LessonTypeLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}

// Values _
func (g LessonType) Values() []string {
	return LessonTypeValues()
}

type LessonOperateAction string

const (
	// 操作类型：创建
	LessonOperateActionCreate LessonOperateAction = "create"
	// 操作类型：编辑此日程
	LessonOperateActionEditOne LessonOperateAction = "edit_one"
	// 操作类型：删除此日程
	LessonOperateActionDeleteOne LessonOperateAction = "delete_one"
	// 操作类型：编辑此日程及后续日程
	LessonOperateActionEditAndAfter LessonOperateAction = "edit_after"
	// 操作类型：删除此日程及后续日程
	LessonOperateActionDeleteAndAfter LessonOperateAction = "delete_after"
)

// Values 操作类型列表
func LessonOperateActionValues() []string {
	return []string{
		string(LessonOperateActionCreate),
		string(LessonOperateActionEditOne),
		string(LessonOperateActionDeleteOne),
		string(LessonOperateActionEditAndAfter),
		string(LessonOperateActionDeleteAndAfter),
	}
}

// LessonOperateActionLabels _
func LessonOperateActionLabels() map[string]string {
	return map[string]string{
		"create":       "create",
		"edit_one":     "edit_one",
		"delete_one":   "delete_one",
		"edit_after":   "edit_after",
		"delete_after": "delete_after",
	}
}

// 多语言支持
var lessonOperateActionLangMap = map[string]map[string]string{
	"zh": {
		"create":       "创建",
		"edit_one":     "编辑此日程",
		"delete_one":   "删除此日程",
		"edit_after":   "编辑此日程及后续日程",
		"delete_after": "删除此日程及后续日程",
	},
	"ru": {
		"create":       "создать",
		"edit_one":     "редактировать экземпляр",
		"delete_one":   "удалить экземпляр",
		"edit_after":   "редактировать экземпляр и последующие экземпляры",
		"delete_after": "удалить экземпляр и последующие экземпляры",
	},
}

func LessonOperateActionLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(lessonOperateActionLangMap, langs)
}

// String _
func (g LessonOperateAction) String() string {
	return string(g)
}

// Label _
func (g LessonOperateAction) Label() string {
	return LessonOperateActionLabels()[g.String()]
}

// Label _
func (g LessonOperateAction) LangLabel(langs ...string) (string, bool) {
	labels, b := LessonOperateActionLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}

// Values _
func (g LessonOperateAction) Values() []string {
	return LessonOperateActionValues()
}

type LessonInstanceStatus string

const (
	// 实例状态：未开始
	LessonInstanceStatusNotStarted LessonInstanceStatus = "not_started"
	// 实例状态：上课中
	LessonInstanceStatusInCourse LessonInstanceStatus = "in_course"
	// 实例状态：已结束
	LessonInstanceStatusEnded LessonInstanceStatus = "ended"
	// 实例状态：已过期（未上课但已过上课时间）
	LessonInstanceStatusExpired LessonInstanceStatus = "expired"
	// 实例状态：已准备好
	LessonInstanceStatusReady LessonInstanceStatus = "ready"
)

// Values 实例状态列表
func LessonInstanceStatusValues() []string {
	return []string{
		string(LessonInstanceStatusNotStarted),
		string(LessonInstanceStatusInCourse),
		string(LessonInstanceStatusEnded),
		string(LessonInstanceStatusExpired),
		string(LessonInstanceStatusReady),
	}
}

// LessonInstanceStatusLabels _
func LessonInstanceStatusLabels() map[string]string {
	return map[string]string{
		"not_started": "not_started",
		"in_course":   "in_course",
		"ended":       "ended",
		"expired":     "expired",
		"ready":       "ready",
	}
}

// 多语言支持
var lessonInstanceStatusLangMap = map[string]map[string]string{
	"zh": {
		"not_started": "未开始",
		"in_course":   "上课中",
		"ended":       "已结束",
		"expired":     "已过期",
		"ready":       "已准备好",
	},
	"ru": {
		"not_started": "не начато",
		"in_course":   "в прямом эфире",
		"ended":       "завершено",
		"expired":     "истекло",
		"ready":       "готов",
	},
}

func LessonInstanceStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(lessonInstanceStatusLangMap, langs)
}

// Label _
func (g LessonInstanceStatus) Label() string {
	return LessonInstanceStatusLabels()[g.String()]
}

// Label _
func (g LessonInstanceStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := LessonInstanceStatusLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}

// Values _
func (g LessonInstanceStatus) Values() []string {
	return LessonInstanceStatusValues()
}

// String _
func (g LessonInstanceStatus) String() string {
	return string(g)
}

type LessonClassKind string

const (
	// 课堂类型：名师课堂
	LessonClassKindFamousTeacher LessonClassKind = "famous_teacher"
	// 课堂类型：名校课堂
	LessonClassKindFamousSchool LessonClassKind = "famous_school"
	// 课堂类型：专递课堂
	LessonClassKindSpecial LessonClassKind = "special"
	// 课堂类型：其他
	LessonClassKindOther LessonClassKind = "other"
)

// Values 课堂类型列表
func LessonClassKindValues() []string {
	return []string{
		string(LessonClassKindFamousTeacher),
		string(LessonClassKindFamousSchool),
		string(LessonClassKindSpecial),
		string(LessonClassKindOther),
	}
}

// LessonClassKindLabels _
func LessonClassKindLabels() map[string]string {
	return map[string]string{
		"famous_teacher": "famous_teacher",
		"famous_school":  "famous_school",
		"special":        "special",
		"other":          "other",
	}
}

// 多语言支持
var lessonClassKindLangMap = map[string]map[string]string{
	"zh": {
		"famous_teacher": "名师课堂",
		"famous_school":  "名校课堂",
		"special":        "专递课堂",
		"other":          "其他",
	},
	"ru": {
		"famous_teacher": "учитель",
		"famous_school":  "школа",
		"special":        "специальный",
		"other":          "другой",
	},
}

// LessonClassKindLangLabels _
func LessonClassKindLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(lessonClassKindLangMap, langs)
}

// String _
func (g LessonClassKind) String() string {
	return string(g)
}

// Label _
func (g LessonClassKind) Label() string {
	return LessonClassKindLabels()[g.String()]
}

// Label _
func (g LessonClassKind) LangLabel(langs ...string) (string, bool) {
	labels, b := LessonClassKindLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}

// Values _
func (g LessonClassKind) Values() []string {
	return LessonClassKindValues()
}

type CourseRecordKind string

const (
	// 课程录制类型：教学课堂
	CourseRecordKindTeachLesson CourseRecordKind = "teach_lesson"
	// 课程录制类型：教研培训
	CourseRecordKindResearchTraining CourseRecordKind = "research_training"
	// 课程录制类型：其他
	CourseRecordKindOther CourseRecordKind = "other"
)

// Values 课程录制类型列表
func CourseRecordKindValues() []string {
	return []string{
		string(CourseRecordKindTeachLesson),
		string(CourseRecordKindResearchTraining),
		string(CourseRecordKindOther),
	}
}

// CourseRecordKindLabels _
func CourseRecordKindLabels() map[string]string {
	return map[string]string{
		"teach_lesson":      "teach_lesson",
		"research_training": "research_training",
		"other":             "other",
	}
}

// 多语言支持
var courseRecordKindLangMap = map[string]map[string]string{
	"zh": {
		"teach_lesson":      "教学课堂",
		"research_training": "教研培训",
		"other":             "其他",
	},
	"ru": {
		"teach_lesson":      "урок",
		"research_training": "исследование",
		"other":             "другой",
	},
}

// CourseRecordKindLangLabels _
func CourseRecordKindLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(courseRecordKindLangMap, langs)
}

// String _
func (g CourseRecordKind) String() string {
	return string(g)
}

// Label _
func (g CourseRecordKind) Label() string {
	return CourseRecordKindLabels()[g.String()]
}

// Label _
func (g CourseRecordKind) LangLabel(langs ...string) (string, bool) {
	labels, b := CourseRecordKindLangLabels(langs...)
	if !b {
		return g.Label(), false
	}
	return labels[g.String()], true
}

// Values _
func (g CourseRecordKind) Values() []string {
	return CourseRecordKindValues()
}
