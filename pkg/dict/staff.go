package dict

import "github.com/zhixingdx/public/pkg/lang"

// EmployeeStatus 员工状态
type EmployeeStatus string

const (
	// 在职
	EmployeeStatusActive EmployeeStatus = "active"
	// 离职
	EmployeeStatusTerminated EmployeeStatus = "terminated"
	// 休假
	EmployeeStatusOnLeave EmployeeStatus = "onLeave"
	// 待入职
	EmployeeStatusPending EmployeeStatus = "pending"
)

// EmployeeStatusValues 员工状态列表
func (s EmployeeStatus) EmployeeStatusValues() []string {
	return []string{
		string(EmployeeStatusActive),
		string(EmployeeStatusTerminated),
		string(EmployeeStatusOnLeave),
		string(EmployeeStatusPending),
	}
}

// EmployeeStatusLabels 员工状态标签
func EmployeeStatusLabels() map[string]string {
	return map[string]string{
		"active":     "在职",
		"terminated": "离职",
		"onLeave":    "休假",
		"pending":    "待入职",
	}
}

// EmployeeStatusLangLabels 员工状态多语言标签
var employeeStatusLangMap = map[string]map[string]string{
	"zh": {
		"active":     "在职",
		"terminated": "离职",
		"onLeave":    "休假",
		"pending":    "待入职",
	},
	"ru": {
		"active":     "работает",
		"terminated": "уволен",
		"onLeave":    "в отпуске",
		"pending":    "ожидает трудоустройства",
	},
}

// EmployeeStatusLangLabels 员工状态多语言标签
func EmployeeStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(employeeStatusLangMap, langs)
}

// String _
func (s EmployeeStatus) String() string {
	return string(s)
}

// Label _
func (s EmployeeStatus) Label() string {
	return EmployeeStatusLabels()[s.String()]
}

// LangLabel _
func (s EmployeeStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := EmployeeStatusLangLabels(langs...)
	if !b {
		return s.Label(), false
	}
	return labels[s.String()], true
}

// Values _
func (s EmployeeStatus) Values() []string {
	return s.EmployeeStatusValues()
}
