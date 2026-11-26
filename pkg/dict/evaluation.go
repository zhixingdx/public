package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

// CourseEvalInviteStatus 评课邀请状态
type CourseEvalInviteStatus string

const (
	// CourseEvalInviteStatusPending 待处理 邀请已发送，等待被邀请人处理
	CourseEvalInviteStatusPending CourseEvalInviteStatus = "pending"
	// CourseEvalInviteStatusAccepted 已接受 被邀请人已接受邀请
	CourseEvalInviteStatusAccepted CourseEvalInviteStatus = "accepted"
	// CourseEvalInviteStatusRejected 已拒绝 被邀请人已拒绝邀请
	CourseEvalInviteStatusRejected CourseEvalInviteStatus = "rejected"
	// CourseEvalInviteStatusExpired 已过期 邀请在规定时间内未处理，自动失效
	CourseEvalInviteStatusExpired CourseEvalInviteStatus = "expired"
	// CourseEvalInviteStatusCancelled 已取消 邀请被取消，不再有效
	CourseEvalInviteStatusCancelled CourseEvalInviteStatus = "cancelled"
	// CourseEvalInviteStatusCompleted 已完成 评课已完成
	CourseEvalInviteStatusCompleted CourseEvalInviteStatus = "completed"
)

// CourseEvalInviteStatusValues 返回所有评课邀请状态的值
func CourseEvalInviteStatusValues() []string {
	return []string{
		string(CourseEvalInviteStatusPending),
		string(CourseEvalInviteStatusAccepted),
		string(CourseEvalInviteStatusRejected),
		string(CourseEvalInviteStatusExpired),
		string(CourseEvalInviteStatusCancelled),
		string(CourseEvalInviteStatusCompleted),
	}
}

// CourseEvalInviteStatusLabels 返回所有评课邀请状态的标签
func CourseEvalInviteStatusLabels() map[string]string {
	return map[string]string{
		"pending":   "pending",
		"accepted":  "accepted",
		"rejected":  "rejected",
		"expired":   "expired",
		"cancelled": "cancelled",
		"completed": "completed",
	}
}

// 多语言支持
var courseEvalInviteStatusLangMap = map[string]map[string]string{
	"zh": {
		"pending":   "待处理",
		"accepted":  "已接受",
		"rejected":  "已拒绝",
		"expired":   "已过期",
		"cancelled": "已取消",
		"completed": "已完成",
	},
	"ru": {
		"pending":   "ожидание",
		"accepted":  "принято",
		"rejected":  "отклонено",
		"expired":   "истек",
		"cancelled": "отменено",
		"completed": "завершено",
	},
}

// CourseEvalInviteStatusLangLabels 返回指定语言的评课邀请状态标签
func CourseEvalInviteStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(courseEvalInviteStatusLangMap, langs)
}

// String 返回字符串
func (c CourseEvalInviteStatus) String() string {
	return string(c)
}

// Label 返回标签
func (c CourseEvalInviteStatus) Label() string {
	return CourseEvalInviteStatusLabels()[c.String()]
}

// LangLabel 返回多语言标签
func (c CourseEvalInviteStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := CourseEvalInviteStatusLangLabels(langs...)
	if !b {
		return c.Label(), false
	}
	return labels[c.String()], true
}

// Values 返回所有值
func (c CourseEvalInviteStatus) Values() []string {
	return CourseEvalInviteStatusValues()
}

// AiAnalysisCompareType 智课分析对比类型
type AiAnalysisCompareType string

const (
	// AiAnalysisCompareTypeSameCourse 同课异构
	AiAnalysisCompareTypeSameCourse AiAnalysisCompareType = "sameCourse"
	// AiAnalysisCompareTypeQualityCourse 好课对比
	AiAnalysisCompareTypeQualityCourse AiAnalysisCompareType = "qualityCourse"
	// AiAnalysisCompareTypePolishCourse 磨课
	AiAnalysisCompareTypePolishCourse AiAnalysisCompareType = "polishCourse"
)

// AiAnalysisCompareTypeValues 返回所有智课分析对比类型
func AiAnalysisCompareTypeValues() []string {
	return []string{
		string(AiAnalysisCompareTypeSameCourse),
		string(AiAnalysisCompareTypeQualityCourse),
		string(AiAnalysisCompareTypePolishCourse),
	}
}

// AiAnalysisCompareTypeLabels 返回所有智课分析对比类型标签
func AiAnalysisCompareTypeLabels() map[string]string {
	return map[string]string{
		"sameCourse":    "同课异构",
		"qualityCourse": "好课对比",
		"polishCourse":  "磨课",
	}
}

// AiAnalysisCompareTypeLangLabels 返回指定语言的智课分析对比类型标签
func AiAnalysisCompareTypeLangLabels(langs ...string) (map[string]string, bool) {
	langMap := map[string]map[string]string{
		"zh": {
			"sameCourse":    "同课异构",
			"qualityCourse": "好课对比",
			"polishCourse":  "磨课",
		},
		"en": {
			"sameCourse":    "Same Course",
			"qualityCourse": "Quality Course",
			"polishCourse":  "Polish Course",
		},
		"ru": {
			"sameCourse":    "сравнение одного курса",
			"qualityCourse": "сравнение качественных уроков",
			"polishCourse":  "шлифовка урока",
		},
	}
	return lang.MatchContent(langMap, langs)
}

// String 返回字符串
func (a AiAnalysisCompareType) String() string {
	return string(a)
}

// Label 返回标签
func (a AiAnalysisCompareType) Label() string {
	return AiAnalysisCompareTypeLabels()[a.String()]
}

// LangLabel 返回多语言标签
func (a AiAnalysisCompareType) LangLabel(langs ...string) (string, bool) {
	labels, b := AiAnalysisCompareTypeLangLabels(langs...)
	if !b {
		return a.Label(), false
	}
	return labels[a.String()], true
}

// Values 返回所有值
func (a AiAnalysisCompareType) Values() []string {
	return AiAnalysisCompareTypeValues()
}

// AiAnalysisCompareDimension 智课分析对比指标维度
type AiAnalysisCompareDimension string

const (
	// AiAnalysisCompareDimensionContent 内容分析
	AiAnalysisCompareDimensionContent AiAnalysisCompareDimension = "content"
	// AiAnalysisCompareDimensionBehavior 行为分析
	AiAnalysisCompareDimensionBehavior AiAnalysisCompareDimension = "behavior"
	// AiAnalysisCompareDimensionThinking 思维能力分析
	AiAnalysisCompareDimensionThinking AiAnalysisCompareDimension = "thinking"
)

// AiAnalysisCompareDimensionValues 返回所有智课分析对比指标维度
func AiAnalysisCompareDimensionValues() []string {
	return []string{
		string(AiAnalysisCompareDimensionContent),
		string(AiAnalysisCompareDimensionBehavior),
		string(AiAnalysisCompareDimensionThinking),
	}
}

// AiAnalysisCompareDimensionLabels 返回所有智课分析对比指标维度标签
func AiAnalysisCompareDimensionLabels() map[string]string {
	return map[string]string{
		"content":  "内容分析",
		"behavior": "行为分析",
		"thinking": "思维能力分析",
	}
}
