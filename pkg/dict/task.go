// cSpell:disable
// Package dict Task
package dict

import "github.com/zhixingdx/public/pkg/lang"

// TaskStatus 用于标识任务状态
type TaskStatus string

const (
	// TaskStatusPending 待分配 任务已创建，等待分配给执行人
	TaskStatusPending TaskStatus = "pending"
	// TaskStatusAssigned 已分配 任务已分配给执行人，等待开始执行
	TaskStatusAssigned TaskStatus = "assigned"
	// TaskStatusInProgress 进行中 执行人已开始执行任务
	TaskStatusInProgress TaskStatus = "inProgress"
	// TaskStatusCompleted 已完成 任务已成功完成
	TaskStatusCompleted TaskStatus = "completed"
	// TaskStatusCancelled 已取消 任务被取消，不再执行
	TaskStatusCancelled TaskStatus = "cancelled"
	// TaskStatusOverdue 已超时 任务超过截止时间未完成
	TaskStatusOverdue TaskStatus = "overdue"
	// TaskStatusFailed 执行失败 任务执行过程中发生错误
	TaskStatusFailed TaskStatus = "failed"
)

// TaskStatusValues 返回所有任务状态的值
func TaskStatusValues() []string {
	return []string{
		string(TaskStatusPending),
		string(TaskStatusAssigned),
		string(TaskStatusInProgress),
		string(TaskStatusCompleted),
		string(TaskStatusCancelled),
		string(TaskStatusOverdue),
		string(TaskStatusFailed),
	}
}

// TaskStatusLabels 返回所有任务状态的标签
func TaskStatusLabels() map[string]string {
	return map[string]string{
		"pending":    "pending",
		"assigned":   "assigned",
		"inProgress": "inProgress",
		"completed":  "completed",
		"cancelled":  "cancelled",
		"overdue":    "overdue",
		"failed":     "failed",
	}
}

// 多语言支持
var taskStatusLangMap = map[string]map[string]string{
	"zh": {
		"pending":    "待分配",
		"assigned":   "已分配",
		"inProgress": "进行中",
		"completed":  "已完成",
		"cancelled":  "已取消",
		"overdue":    "已超时",
		"failed":     "执行失败",
	},
	"en": {
		"pending":    "Pending",
		"assigned":   "Assigned",
		"inProgress": "In Progress",
		"completed":  "Completed",
		"cancelled":  "Cancelled",
		"overdue":    "Overdue",
		"failed":     "Failed",
	},
	"ru": {
		"pending":    "ожидание",
		"assigned":   "назначенный",
		"inProgress": "в процессе",
		"completed":  "завершенный",
		"cancelled":  "отмененный",
		"overdue":    "просроченный",
		"failed":     "неудачный",
	},
}

// TaskStatusLangLabels 返回指定语言的任务状态标签
func TaskStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(taskStatusLangMap, langs)
}

// String 返回字符串
func (t TaskStatus) String() string {
	return string(t)
}

// Label 返回标签
func (t TaskStatus) Label() string {
	return TaskStatusLabels()[t.String()]
}

// LangLabel 返回多语言标签
func (t TaskStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := TaskStatusLangLabels(langs...)
	if !b {
		return t.Label(), false
	}
	return labels[t.String()], true
}

// Values 返回所有值
func (t TaskStatus) Values() []string {
	return TaskStatusValues()
}

// IsFinalState 判断是否为最终状态（不可再变更）
func (t TaskStatus) IsFinalState() bool {
	switch t {
	case TaskStatusCompleted, TaskStatusCancelled, TaskStatusFailed:
		return true
	default:
		return false
	}
}

// IsActiveState 判断是否为活跃状态（正在执行中）
func (t TaskStatus) IsActiveState() bool {
	switch t {
	case TaskStatusAssigned, TaskStatusInProgress:
		return true
	default:
		return false
	}
}

// TaskAction 任务操作类型
type TaskAction string

const (
	// TaskActionAssign 分配任务 将任务分配给执行人
	TaskActionAssign TaskAction = "assign"
	// TaskActionStart 开始执行 执行人开始执行任务
	TaskActionStart TaskAction = "start"
	// TaskActionComplete 完成任务 执行人完成任务
	TaskActionComplete TaskAction = "complete"
	// TaskActionCancel 取消任务 取消任务执行
	TaskActionCancel TaskAction = "cancel"
	// TaskActionReassign 重新分配 将任务重新分配给其他执行人
	TaskActionReassign TaskAction = "reassign"
	// TaskActionUpdate 更新进度 更新任务执行进度
	TaskActionUpdate TaskAction = "update"
	// TaskActionPause 暂停任务 暂停任务执行
	TaskActionPause TaskAction = "pause"
	// TaskActionResume 恢复任务 恢复被暂停的任务
	TaskActionResume TaskAction = "resume"
	// TaskActionExpire 过期 任务超过截止时间自动过期
	TaskActionExpire TaskAction = "expire"
	// TaskActionFail 标记失败 标记任务执行失败
	TaskActionFail TaskAction = "fail"
)

// TaskActionValues 返回所有任务操作类型的值
func TaskActionValues() []string {
	return []string{
		string(TaskActionAssign),
		string(TaskActionStart),
		string(TaskActionComplete),
		string(TaskActionCancel),
		string(TaskActionReassign),
		string(TaskActionUpdate),
		string(TaskActionPause),
		string(TaskActionResume),
		string(TaskActionExpire),
		string(TaskActionFail),
	}
}

// TaskActionLabels 返回所有任务操作类型的标签
func TaskActionLabels() map[string]string {
	return map[string]string{
		"assign":   "分配任务",
		"start":    "开始执行",
		"complete": "完成任务",
		"cancel":   "取消任务",
		"reassign": "重新分配",
		"update":   "更新进度",
		"pause":    "暂停任务",
		"resume":   "恢复任务",
		"expire":   "过期",
		"fail":     "标记失败",
	}
}

// 任务操作类型多语言支持
var taskActionLangMap = map[string]map[string]string{
	"zh": {
		"assign":   "分配任务",
		"start":    "开始执行",
		"complete": "完成任务",
		"cancel":   "取消任务",
		"reassign": "重新分配",
		"update":   "更新进度",
		"pause":    "暂停任务",
		"resume":   "恢复任务",
		"expire":   "过期",
		"fail":     "标记失败",
	},
	"en": {
		"assign":   "Assign Task",
		"start":    "Start",
		"complete": "Complete",
		"cancel":   "Cancel",
		"reassign": "Reassign",
		"update":   "Update Progress",
		"pause":    "Pause",
		"resume":   "Resume",
		"expire":   "Expire",
		"fail":     "Mark as Failed",
	},
	"ru": {
		"assign":   "назначить задачу",
		"start":    "начать",
		"complete": "завершить",
		"cancel":   "отменить",
		"reassign": "переназначить",
		"update":   "обновить прогресс",
		"pause":    "приостановить",
		"resume":   "возобновить",
		"expire":   "истечь",
		"fail":     "отметить как неудачный",
	},
}

// TaskActionLangLabels 返回指定语言的任务操作类型标签
func TaskActionLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(taskActionLangMap, langs)
}

// String 返回字符串
func (t TaskAction) String() string {
	return string(t)
}

// Label 返回标签
func (t TaskAction) Label() string {
	return TaskActionLabels()[t.String()]
}

// LangLabel 返回多语言标签
func (t TaskAction) LangLabel(langs ...string) (string, bool) {
	labels, b := TaskActionLangLabels(langs...)
	if !b {
		return t.Label(), false
	}
	return labels[t.String()], true
}

// Values 返回所有值
func (t TaskAction) Values() []string {
	return TaskActionValues()
}

// IsUserOperation 判断是否为用户操作（非系统自动操作）
func (t TaskAction) IsUserOperation() bool {
	switch t {
	case TaskActionExpire:
		return false
	default:
		return true
	}
}

// IsSystemOperation 判断是否为系统操作
func (t TaskAction) IsSystemOperation() bool {
	return !t.IsUserOperation()
}

// IsExecutionOperation 判断是否为执行相关操作
func (t TaskAction) IsExecutionOperation() bool {
	switch t {
	case TaskActionStart, TaskActionComplete, TaskActionUpdate, TaskActionPause, TaskActionResume, TaskActionFail:
		return true
	default:
		return false
	}
}

// IsFlowControlOperation 判断是否为流程控制操作
func (t TaskAction) IsFlowControlOperation() bool {
	switch t {
	case TaskActionAssign, TaskActionCancel, TaskActionReassign, TaskActionExpire:
		return true
	default:
		return false
	}
}

// IsManagementOperation 判断是否为管理操作
func (t TaskAction) IsManagementOperation() bool {
	switch t {
	case TaskActionReassign:
		return true
	default:
		return false
	}
}

// TaskExecutionMode 任务执行模式
type TaskExecutionMode string

const (
	// TaskExecutionModeEmbedded 内嵌式执行 任务服务提供完整的处理界面和逻辑
	TaskExecutionModeEmbedded TaskExecutionMode = "embedded"
	// TaskExecutionModeRedirect 跳转式执行 任务服务只记录状态，实际处理在业务服务中完成
	TaskExecutionModeRedirect TaskExecutionMode = "redirect"
)

// TaskExecutionModeValues 返回所有执行模式的值
func TaskExecutionModeValues() []string {
	return []string{
		string(TaskExecutionModeEmbedded),
		string(TaskExecutionModeRedirect),
	}
}

// TaskExecutionModeLabels 返回所有执行模式的标签
func TaskExecutionModeLabels() map[string]string {
	return map[string]string{
		"embedded": "内嵌式",
		"redirect": "跳转式",
	}
}

// 执行模式多语言支持
var taskExecutionModeLangMap = map[string]map[string]string{
	"zh": {
		"embedded": "内嵌式",
		"redirect": "跳转式",
	},
	"en": {
		"embedded": "Embedded",
		"redirect": "Redirect",
	},
	"ru": {
		"embedded": "встроенный",
		"redirect": "перенаправление",
	},
}

// TaskExecutionModeLangLabels 返回指定语言的执行模式标签
func TaskExecutionModeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(taskExecutionModeLangMap, langs)
}

// String 返回字符串
func (t TaskExecutionMode) String() string {
	return string(t)
}

// Label 返回标签
func (t TaskExecutionMode) Label() string {
	return TaskExecutionModeLabels()[t.String()]
}

// LangLabel 返回多语言标签
func (t TaskExecutionMode) LangLabel(langs ...string) (string, bool) {
	labels, b := TaskExecutionModeLangLabels(langs...)
	if !b {
		return t.Label(), false
	}
	return labels[t.String()], true
}

// Values 返回所有值
func (t TaskExecutionMode) Values() []string {
	return TaskExecutionModeValues()
}

// TaskRecordKind 任务记录类型（操作来源）
type TaskRecordKind string

const (
	// TaskRecordKindUser 用户操作 用户直接操作任务
	TaskRecordKindUser TaskRecordKind = "user"
	// TaskRecordKindSystem 系统操作 系统自动操作（如过期、超时等）
	TaskRecordKindSystem TaskRecordKind = "system"
	// TaskRecordKindBusiness 业务服务回调 跳转式任务由业务服务回调更新状态
	TaskRecordKindBusiness TaskRecordKind = "business"
)

// TaskRecordKindValues 返回所有任务记录类型的值
func TaskRecordKindValues() []string {
	return []string{
		string(TaskRecordKindUser),
		string(TaskRecordKindSystem),
		string(TaskRecordKindBusiness),
	}
}

// TaskRecordKindLabels 返回所有任务记录类型的标签
func TaskRecordKindLabels() map[string]string {
	return map[string]string{
		"user":     "用户操作",
		"system":   "系统操作",
		"business": "业务服务回调",
	}
}

// 任务记录类型多语言支持
var taskRecordKindLangMap = map[string]map[string]string{
	"zh": {
		"user":     "用户操作",
		"system":   "系统操作",
		"business": "业务服务回调",
	},
	"en": {
		"user":     "User Operation",
		"system":   "System Operation",
		"business": "Business Callback",
	},
	"ru": {
		"user":     "действие пользователя",
		"system":   "системное действие",
		"business": "бизнес-обратный вызов",
	},
}

// TaskRecordKindLangLabels 返回指定语言的任务记录类型标签
func TaskRecordKindLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(taskRecordKindLangMap, langs)
}

// String 返回字符串
func (t TaskRecordKind) String() string {
	return string(t)
}

// Label 返回标签
func (t TaskRecordKind) Label() string {
	return TaskRecordKindLabels()[t.String()]
}

// LangLabel 返回多语言标签
func (t TaskRecordKind) LangLabel(langs ...string) (string, bool) {
	labels, b := TaskRecordKindLangLabels(langs...)
	if !b {
		return t.Label(), false
	}
	return labels[t.String()], true
}

// Values 返回所有值
func (t TaskRecordKind) Values() []string {
	return TaskRecordKindValues()
}
