// cSpell:disable
// Package dict Workflow
package dict

import "github.com/zhixingdx/oipivot/backend/pkg/lang"

// WorkflowInstanceStatus 工作流实例状态
type WorkflowInstanceStatus string

const (
	// WorkflowInstanceStatusPending 等待执行 实例已创建但未开始
	WorkflowInstanceStatusPending WorkflowInstanceStatus = "pending"
	// WorkflowInstanceStatusRunning 执行中 工作流正在运行
	WorkflowInstanceStatusRunning WorkflowInstanceStatus = "running"
	// WorkflowInstanceStatusPaused 已暂停 等待恢复或人工干预
	WorkflowInstanceStatusPaused WorkflowInstanceStatus = "paused"
	// WorkflowInstanceStatusCompleted 已完成 执行成功
	WorkflowInstanceStatusCompleted WorkflowInstanceStatus = "completed"
	// WorkflowInstanceStatusFailed 已失败 执行出错
	WorkflowInstanceStatusFailed WorkflowInstanceStatus = "failed"
	// WorkflowInstanceStatusCancelled 已取消 被主动终止
	WorkflowInstanceStatusCancelled WorkflowInstanceStatus = "cancelled"
	// WorkflowInstanceStatusTimeout 已超时 执行时间超过限制
	WorkflowInstanceStatusTimeout WorkflowInstanceStatus = "timeout"
)

// WorkflowInstanceStatusValues 返回所有工作流实例状态的值
func WorkflowInstanceStatusValues() []string {
	return []string{
		string(WorkflowInstanceStatusPending),
		string(WorkflowInstanceStatusRunning),
		string(WorkflowInstanceStatusPaused),
		string(WorkflowInstanceStatusCompleted),
		string(WorkflowInstanceStatusFailed),
		string(WorkflowInstanceStatusCancelled),
		string(WorkflowInstanceStatusTimeout),
	}
}

// WorkflowInstanceStatusLabels 返回所有工作流实例状态的标签
func WorkflowInstanceStatusLabels() map[string]string {
	return map[string]string{
		"pending":   "pending",
		"running":   "running",
		"paused":    "paused",
		"completed": "completed",
		"failed":    "failed",
		"cancelled": "cancelled",
		"timeout":   "timeout",
	}
}

// 多语言支持
var workflowInstanceStatusLangMap = map[string]map[string]string{
	"zh": {
		"pending":   "等待执行",
		"running":   "执行中",
		"paused":    "已暂停",
		"completed": "已完成",
		"failed":    "已失败",
		"cancelled": "已取消",
		"timeout":   "已超时",
	},
	"en": {
		"pending":   "Pending",
		"running":   "Running",
		"paused":    "Paused",
		"completed": "Completed",
		"failed":    "Failed",
		"cancelled": "Cancelled",
		"timeout":   "Timeout",
	},
	"ru": {
		"pending":   "ожидание",
		"running":   "выполняется",
		"paused":    "приостановлено",
		"completed": "завершено",
		"failed":    "неудача",
		"cancelled": "отменено",
		"timeout":   "тайм-аут",
	},
}

// WorkflowInstanceStatusLangLabels 返回指定语言的工作流实例状态标签
func WorkflowInstanceStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(workflowInstanceStatusLangMap, langs)
}

// String 返回字符串
func (w WorkflowInstanceStatus) String() string {
	return string(w)
}

// Label 返回标签
func (w WorkflowInstanceStatus) Label() string {
	return WorkflowInstanceStatusLabels()[w.String()]
}

// LangLabel 返回多语言标签
func (w WorkflowInstanceStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := WorkflowInstanceStatusLangLabels(langs...)
	if !b {
		return w.Label(), false
	}
	return labels[w.String()], true
}

// Values 返回所有值
func (w WorkflowInstanceStatus) Values() []string {
	return WorkflowInstanceStatusValues()
}

// IsFinalState 判断是否为最终状态（不可再变更）
func (w WorkflowInstanceStatus) IsFinalState() bool {
	switch w {
	case WorkflowInstanceStatusCompleted, WorkflowInstanceStatusFailed,
		WorkflowInstanceStatusCancelled, WorkflowInstanceStatusTimeout:
		return true
	default:
		return false
	}
}

// IsActiveState 判断是否为活跃状态（正在执行中）
func (w WorkflowInstanceStatus) IsActiveState() bool {
	switch w {
	case WorkflowInstanceStatusRunning, WorkflowInstanceStatusPaused:
		return true
	default:
		return false
	}
}

// WorkflowNodeStatus 工作流节点状态
type WorkflowNodeStatus string

const (
	// WorkflowNodeStatusPending 等待执行
	WorkflowNodeStatusPending WorkflowNodeStatus = "pending"
	// WorkflowNodeStatusRunning 执行中
	WorkflowNodeStatusRunning WorkflowNodeStatus = "running"
	// WorkflowNodeStatusCompleted 执行成功
	WorkflowNodeStatusCompleted WorkflowNodeStatus = "completed"
	// WorkflowNodeStatusFailed 执行失败
	WorkflowNodeStatusFailed WorkflowNodeStatus = "failed"
	// WorkflowNodeStatusSkipped 已跳过（条件不满足）
	WorkflowNodeStatusSkipped WorkflowNodeStatus = "skipped"
	// WorkflowNodeStatusCancelled 已取消
	WorkflowNodeStatusCancelled WorkflowNodeStatus = "cancelled"
)

// WorkflowNodeStatusValues 返回所有节点状态的值
func WorkflowNodeStatusValues() []string {
	return []string{
		string(WorkflowNodeStatusPending),
		string(WorkflowNodeStatusRunning),
		string(WorkflowNodeStatusCompleted),
		string(WorkflowNodeStatusFailed),
		string(WorkflowNodeStatusSkipped),
		string(WorkflowNodeStatusCancelled),
	}
}

// WorkflowNodeStatusLabels 返回所有节点状态的标签
func WorkflowNodeStatusLabels() map[string]string {
	return map[string]string{
		"pending":   "pending",
		"running":   "running",
		"completed": "completed",
		"failed":    "failed",
		"skipped":   "skipped",
		"cancelled": "cancelled",
	}
}

// 节点状态多语言支持
var workflowNodeStatusLangMap = map[string]map[string]string{
	"zh": {
		"pending":   "等待执行",
		"running":   "执行中",
		"completed": "已完成",
		"failed":    "已失败",
		"skipped":   "已跳过",
		"cancelled": "已取消",
	},
	"en": {
		"pending":   "Pending",
		"running":   "Running",
		"completed": "Completed",
		"failed":    "Failed",
		"skipped":   "Skipped",
		"cancelled": "Cancelled",
	},
	"ru": {
		"pending":   "ожидание",
		"running":   "выполняется",
		"completed": "завершено",
		"failed":    "неудача",
		"skipped":   "пропущено",
		"cancelled": "отменено",
	},
}

// WorkflowNodeStatusLangLabels 返回指定语言的节点状态标签
func WorkflowNodeStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(workflowNodeStatusLangMap, langs)
}

// String 返回字符串
func (w WorkflowNodeStatus) String() string {
	return string(w)
}

// Label 返回标签
func (w WorkflowNodeStatus) Label() string {
	return WorkflowNodeStatusLabels()[w.String()]
}

// LangLabel 返回多语言标签
func (w WorkflowNodeStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := WorkflowNodeStatusLangLabels(langs...)
	if !b {
		return w.Label(), false
	}
	return labels[w.String()], true
}

// Values 返回所有值
func (w WorkflowNodeStatus) Values() []string {
	return WorkflowNodeStatusValues()
}

// IsFinalState 判断是否为最终状态
func (w WorkflowNodeStatus) IsFinalState() bool {
	switch w {
	case WorkflowNodeStatusCompleted, WorkflowNodeStatusFailed,
		WorkflowNodeStatusSkipped, WorkflowNodeStatusCancelled:
		return true
	default:
		return false
	}
}

// WorkflowNodeType 工作流节点类型
type WorkflowNodeType string

const (
	// WorkflowNodeTypeCall 调用任务节点
	WorkflowNodeTypeCall WorkflowNodeType = "call"
	// WorkflowNodeTypeFork 并行分支节点
	WorkflowNodeTypeFork WorkflowNodeType = "fork"
	// WorkflowNodeTypeSwitch 条件分支节点
	WorkflowNodeTypeSwitch WorkflowNodeType = "switch"
	// WorkflowNodeTypeFor 循环节点
	WorkflowNodeTypeFor WorkflowNodeType = "for"
	// WorkflowNodeTypeTry 异常处理节点
	WorkflowNodeTypeTry WorkflowNodeType = "try"
	// WorkflowNodeTypeWait 等待节点
	WorkflowNodeTypeWait WorkflowNodeType = "wait"
	// WorkflowNodeTypeSet 变量设置节点
	WorkflowNodeTypeSet WorkflowNodeType = "set"
)

// WorkflowNodeTypeValues 返回所有节点类型的值
func WorkflowNodeTypeValues() []string {
	return []string{
		string(WorkflowNodeTypeCall),
		string(WorkflowNodeTypeFork),
		string(WorkflowNodeTypeSwitch),
		string(WorkflowNodeTypeFor),
		string(WorkflowNodeTypeTry),
		string(WorkflowNodeTypeWait),
		string(WorkflowNodeTypeSet),
	}
}

// WorkflowNodeTypeLabels 返回所有节点类型的标签
func WorkflowNodeTypeLabels() map[string]string {
	return map[string]string{
		"call":   "调用任务",
		"fork":   "并行分支",
		"switch": "条件分支",
		"for":    "循环",
		"try":    "异常处理",
		"wait":   "等待",
		"set":    "变量设置",
	}
}

// 节点类型多语言支持
var workflowNodeTypeLangMap = map[string]map[string]string{
	"zh": {
		"call":   "调用任务",
		"fork":   "并行分支",
		"switch": "条件分支",
		"for":    "循环",
		"try":    "异常处理",
		"wait":   "等待",
		"set":    "变量设置",
	},
	"en": {
		"call":   "Call Task",
		"fork":   "Fork",
		"switch": "Switch",
		"for":    "For Loop",
		"try":    "Try Catch",
		"wait":   "Wait",
		"set":    "Set Variable",
	},
	"ru": {
		"call":   "вызов задачи",
		"fork":   "развилка",
		"switch": "переключатель",
		"for":    "цикл",
		"try":    "обработка исключений",
		"wait":   "ожидание",
		"set":    "установка переменной",
	},
}

// WorkflowNodeTypeLangLabels 返回指定语言的节点类型标签
func WorkflowNodeTypeLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(workflowNodeTypeLangMap, langs)
}

// String 返回字符串
func (w WorkflowNodeType) String() string {
	return string(w)
}

// Label 返回标签
func (w WorkflowNodeType) Label() string {
	return WorkflowNodeTypeLabels()[w.String()]
}

// LangLabel 返回多语言标签
func (w WorkflowNodeType) LangLabel(langs ...string) (string, bool) {
	labels, b := WorkflowNodeTypeLangLabels(langs...)
	if !b {
		return w.Label(), false
	}
	return labels[w.String()], true
}

// Values 返回所有值
func (w WorkflowNodeType) Values() []string {
	return WorkflowNodeTypeValues()
}
