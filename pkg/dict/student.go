package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

// AcademicStatus 学籍状态
type AcademicStatus string

const (
	// 在读
	AcademicStatusEnrolled AcademicStatus = "enrolled"
	// 休学
	AcademicStatusSuspended AcademicStatus = "suspended"
	// 退学
	AcademicStatusWithdrawn AcademicStatus = "withdrawn"
	// 转学
	AcademicStatusTransferred AcademicStatus = "transferredOut"
	// 毕业
	AcademicStatusGraduated AcademicStatus = "graduated"
)

// Values 学籍状态列表
func (s AcademicStatus) AcademicStatusValues() []string {
	return []string{
		string(AcademicStatusEnrolled),
		string(AcademicStatusSuspended),
		string(AcademicStatusWithdrawn),
		string(AcademicStatusTransferred),
		string(AcademicStatusGraduated),
	}
}

// AcademicStatusLabels 学籍状态标签
func AcademicStatusLabels() map[string]string {
	return map[string]string{
		"enrolled":    "在读",
		"suspended":   "休学",
		"withdrawn":   "退学",
		"transferred": "转学",
		"graduated":   "毕业",
	}
}

// AcademicStatusLangLabels 学籍状态多语言标签
var academicStatusLangMap = map[string]map[string]string{
	"zh": {
		"enrolled":    "在读",
		"suspended":   "休学",
		"withdrawn":   "退学",
		"transferred": "转学",
		"graduated":   "毕业",
	},
	"ru": {
		"enrolled":    "вобучении",
		"suspended":   "отчисленный",
		"withdrawn":   "выбывший",
		"transferred": "переводящийся",
		"graduated":   "выпускник",
	},
}

func AcademicStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(academicStatusLangMap, langs)
}

// String _
func (s AcademicStatus) String() string {
	return string(s)
}

// Label _
func (s AcademicStatus) Label() string {
	return AcademicStatusLabels()[s.String()]
}

// LangLabel _
func (s AcademicStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := AcademicStatusLangLabels(langs...)
	if !b {
		return s.Label(), false
	}
	return labels[s.String()], true
}

// Values _
func (s AcademicStatus) Values() []string {
	return s.AcademicStatusValues()
}

// ClassRelationType 学生班级关系类型
type ClassRelationType string

const (
	// 新生入班
	ClassRelationTypeEnrollment ClassRelationType = "enrollment"
	// 转班
	ClassRelationTypeTransfer ClassRelationType = "transfer"
	// 升学升班
	ClassRelationTypePromotion ClassRelationType = "promotion"
	// 留级/复读
	ClassRelationTypeRepeat ClassRelationType = "repeat"
	// 复学入班
	ClassRelationTypeReEnrollment ClassRelationType = "reEnrollment"
	// 暂时离班
	ClassRelationTypeTemporaryDeparture ClassRelationType = "temporaryDeparture"
)

// ClassRelationTypeValues 班级关系类型列表
func (s ClassRelationType) ClassRelationTypeValues() []string {
	return []string{
		string(ClassRelationTypeEnrollment),
		string(ClassRelationTypeTransfer),
		string(ClassRelationTypePromotion),
		string(ClassRelationTypeRepeat),
		string(ClassRelationTypeReEnrollment),
		string(ClassRelationTypeTemporaryDeparture),
	}
}

// ClassRelationTypeLabels 班级关系类型标签
func (s ClassRelationType) ClassRelationTypeLabels() map[string]string {
	return map[string]string{
		"enrollment":         "新生入班",
		"transfer":           "转班",
		"promotion":          "升学升班",
		"repeat":             "留级/复读",
		"reEnrollment":       "复学入班",
		"temporaryDeparture": "暂时离班",
	}
}

// ClassRelationTypeLangMap 班级关系类型多语言标签
var classRelationTypeLangMap = map[string]map[string]string{
	"zh": {
		"enrollment":         "新生入班",
		"transfer":           "转班",
		"promotion":          "升学升班",
		"repeat":             "留级/复读",
		"reEnrollment":       "复学入班",
		"temporaryDeparture": "暂时离班",
	},
	"ru": {
		"enrollment":         "новый класс",
		"transfer":           "перевод",
		"promotion":          "повышение",
		"repeat":             "повтор",
		"reEnrollment":       "повторный класс",
		"temporaryDeparture": "временный отпуск",
	},
}

// ClassRelationTypeLangLabels 班级关系类型多语言标签
func (s ClassRelationType) ClassRelationTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(classRelationTypeLangMap, langs)
}

// String _
func (s ClassRelationType) String() string {
	return string(s)
}

// Label _
func (s ClassRelationType) Label() string {
	return s.ClassRelationTypeLabels()[s.String()]
}

// LangLabel _
func (s ClassRelationType) LangLabel(langs ...string) (string, bool) {
	labels, b := s.ClassRelationTypeLangLabels(langs...)
	if !b {
		return s.ClassRelationTypeLabels()[s.ClassRelationTypeValues()[0]], false
	}
	return labels[s.ClassRelationTypeValues()[0]], true
}

// Values _
func (s ClassRelationType) Values() []string {
	return s.ClassRelationTypeValues()
}

type ClassRelationStatus string

const (
	// 待生效
	ClassRelationStatusPending ClassRelationStatus = "pending"
	// 生效
	ClassRelationStatusActive ClassRelationStatus = "active"
	// 失效
	ClassRelationStatusExpired ClassRelationStatus = "expired"
)

// ClassRelationStatusValues 班级关系状态列表
func (s ClassRelationStatus) ClassRelationStatusValues() []string {
	return []string{
		string(ClassRelationStatusPending),
		string(ClassRelationStatusActive),
		string(ClassRelationStatusExpired),
	}
}

// ClassRelationStatusLabels 班级关系状态标签
func (s ClassRelationStatus) ClassRelationStatusLabels() map[string]string {
	return map[string]string{
		"pending": "待生效",
		"active":  "生效",
		"expired": "失效",
	}
}

// ClassRelationStatusLangMap 班级关系状态多语言标签
var classRelationStatusLangMap = map[string]map[string]string{
	"zh": {
		"pending": "待生效",
		"active":  "生效",
		"expired": "失效",
	},
	"ru": {
		"pending": "ожидание",
		"active":  "активно",
		"expired": "устарело",
	},
}

// ClassRelationStatusLangLabels 班级关系状态多语言标签
func (s ClassRelationStatus) ClassRelationStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(classRelationStatusLangMap, langs)
}

// String _
func (s ClassRelationStatus) String() string {
	return string(s)
}

// Label _
func (s ClassRelationStatus) Label() string {
	return s.ClassRelationStatusLabels()[s.String()]
}

// LangLabel _
func (s ClassRelationStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := s.ClassRelationStatusLangLabels(langs...)
	if !b {
		return s.ClassRelationStatusLabels()[s.ClassRelationStatusValues()[0]], false
	}
	return labels[s.String()], true
}

// Values _
func (s ClassRelationStatus) Values() []string {
	return s.ClassRelationStatusValues()
}
