// cSpell:disable
// Package dict Approval
package dict

import "github.com/zhixingdx/public/pkg/lang"

// ApprovalStatus 用于标识审批状态
type ApprovalStatus string

const (
	// ApprovalStatusPending 待审批 请求提交后，等待审批人的审核。
	ApprovalStatusPending ApprovalStatus = "pending"
	// ApprovalStatusIn 审批中 审批人正在审核请求。
	ApprovalStatusIn ApprovalStatus = "in"
	// ApprovalStatusApproved 已批准 请求已经得到审批人的同意。
	ApprovalStatusApproved ApprovalStatus = "approved"
	// ApprovalStatusRejected 拒绝 请求未通过审批人的审核。
	ApprovalStatusRejected ApprovalStatus = "rejected"
	// ApprovalStatusReturned 退回修改 请求需要修改后重新提交审批。
	ApprovalStatusReturned ApprovalStatus = "returned"
	// ApprovalStatusWithdrawn 撤销 请求者主动撤销了请求。
	ApprovalStatusWithdrawn ApprovalStatus = "withdrawn"
	// ApprovalStatusExpired 过期 请求在规定时间内未得到审批，自动失效。
	ApprovalStatusExpired ApprovalStatus = "expired"
	// ApprovalStatusSuspended 挂起 审批流程因某些原因暂时中断。
	ApprovalStatusSuspended ApprovalStatus = "suspended"
	// ApprovalStatusCompleted 完成 审批流程结束，无论结果如何。
	ApprovalStatusCompleted ApprovalStatus = "completed"
)

// ApprovalStatusValues 返回所有审批状态的值
func ApprovalStatusValues() []string {
	return []string{
		string(ApprovalStatusPending),
		string(ApprovalStatusIn),
		string(ApprovalStatusApproved),
		string(ApprovalStatusRejected),
		string(ApprovalStatusReturned),
		string(ApprovalStatusWithdrawn),
		string(ApprovalStatusExpired),
		string(ApprovalStatusSuspended),
		string(ApprovalStatusCompleted),
	}
}

// ApprovalStatusLabels 返回所有审批状态的标签
func ApprovalStatusLabels() map[string]string {
	return map[string]string{
		"pending":   "pending",
		"in":        "in",
		"approved":  "approved",
		"rejected":  "rejected",
		"returned":  "returned",
		"withdrawn": "withdrawn",
		"expired":   "expired",
		"suspended": "suspended",
		"completed": "completed",
	}
}

// 多语言支持
var approvalStatusLangMap = map[string]map[string]string{
	"zh": {
		"pending":   "待审批",
		"in":        "审批中",
		"approved":  "已批准",
		"rejected":  "拒绝",
		"returned":  "退回修改",
		"withdrawn": "撤销",
		"expired":   "过期",
		"suspended": "挂起",
		"completed": "完成",
	},
	"ru": {
		"pending":   "ожидание",
		"in":        "в процессе",
		"approved":  "утвержденный",
		"rejected":  "отклоненный",
		"returned":  "возвращенный",
		"withdrawn": "отозванный",
		"expired":   "истекший",
		"suspended": "приостановленный",
		"completed": "завершенный",
	},
}

// ApprovalStatusLangLabels 返回指定语言的审批状态标签
func ApprovalStatusLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(approvalStatusLangMap, langs)
}

// String 返回字符串
func (a ApprovalStatus) String() string {
	return string(a)
}

// Label 返回标签
func (a ApprovalStatus) Label() string {
	return ApprovalStatusLabels()[a.String()]
}

// LangLabel
func (a ApprovalStatus) LangLabel(langs ...string) (string, bool) {
	labels, b := ApprovalStatusLangLabels(langs...)
	if !b {
		return a.Label(), false
	}
	return labels[a.String()], true
}

// Values _
func (a ApprovalStatus) Values() []string {
	return ApprovalStatusValues()
}

// ApprovalAction 审批操作类型
type ApprovalAction string

const (
	// ApprovalActionSubmit 提交审批 申请人提交审批请求
	ApprovalActionSubmit ApprovalAction = "submit"
	// ApprovalActionApprove 批准 审批人同意申请
	ApprovalActionApprove ApprovalAction = "approve"
	// ApprovalActionReject 拒绝 审批人拒绝申请
	ApprovalActionReject ApprovalAction = "reject"
	// ApprovalActionReturn 退回修改 审批人要求申请人修改后重新提交
	ApprovalActionReturn ApprovalAction = "return"
	// ApprovalActionWithdraw 撤销 申请人主动撤销申请
	ApprovalActionWithdraw ApprovalAction = "withdraw"
	// ApprovalActionSuspend 挂起 暂停审批流程
	ApprovalActionSuspend ApprovalAction = "suspend"
	// ApprovalActionResume 恢复 恢复被挂起的审批流程
	ApprovalActionResume ApprovalAction = "resume"
	// ApprovalActionTransfer 转交 将审批转交给其他人处理
	ApprovalActionTransfer ApprovalAction = "transfer"
	// ApprovalActionDelegate 委派 委派他人代为审批
	ApprovalActionDelegate ApprovalAction = "delegate"
	// ApprovalActionAddApprover 添加审批人 在审批过程中添加新的审批人
	ApprovalActionAddApprover ApprovalAction = "add_approver"
	// ApprovalActionRemoveApprover 移除审批人 移除某个审批人
	ApprovalActionRemoveApprover ApprovalAction = "remove_approver"
	// ApprovalActionExpire 过期 审批超时自动过期
	ApprovalActionExpire ApprovalAction = "expire"
	// ApprovalActionComplete 完成 审批流程完成
	ApprovalActionComplete ApprovalAction = "complete"
	// ApprovalActionCancel 取消 取消审批流程
	ApprovalActionCancel ApprovalAction = "cancel"
	// ApprovalActionRestart 重启 重新启动审批流程
	ApprovalActionRestart ApprovalAction = "restart"
	// ApprovalActionComment 评论 添加审批评论或意见
	ApprovalActionComment ApprovalAction = "comment"
	// ApprovalActionView 查看 查看审批详情
	ApprovalActionView ApprovalAction = "view"
	// ApprovalActionRemind 提醒 发送审批提醒
	ApprovalActionRemind ApprovalAction = "remind"
	// ApprovalActionEscalate 升级 将审批升级到更高层级
	ApprovalActionEscalate ApprovalAction = "escalate"
	// ApprovalActionAutoApprove 自动批准 系统自动批准
	ApprovalActionAutoApprove ApprovalAction = "auto_approve"
	// ApprovalActionAutoReject 自动拒绝 系统自动拒绝
	ApprovalActionAutoReject ApprovalAction = "auto_reject"
	// ApprovalActionResubmit 再次提交 被拒绝或退回后重新提交审批
	ApprovalActionResubmit ApprovalAction = "resubmit"
)

// ApprovalActionValues 返回所有审批操作类型的值
func ApprovalActionValues() []string {
	return []string{
		string(ApprovalActionSubmit),
		string(ApprovalActionApprove),
		string(ApprovalActionReject),
		string(ApprovalActionReturn),
		string(ApprovalActionWithdraw),
		string(ApprovalActionSuspend),
		string(ApprovalActionResume),
		string(ApprovalActionTransfer),
		string(ApprovalActionDelegate),
		string(ApprovalActionAddApprover),
		string(ApprovalActionRemoveApprover),
		string(ApprovalActionExpire),
		string(ApprovalActionComplete),
		string(ApprovalActionCancel),
		string(ApprovalActionRestart),
		string(ApprovalActionComment),
		string(ApprovalActionView),
		string(ApprovalActionRemind),
		string(ApprovalActionEscalate),
		string(ApprovalActionAutoApprove),
		string(ApprovalActionAutoReject),
		string(ApprovalActionResubmit),
	}
}

// ApprovalActionLabels 返回所有审批操作类型的标签
func ApprovalActionLabels() map[string]string {
	return map[string]string{
		"submit":          "提交审批",
		"approve":         "批准",
		"reject":          "拒绝",
		"return":          "退回修改",
		"withdraw":        "撤销",
		"suspend":         "挂起",
		"resume":          "恢复",
		"transfer":        "转交",
		"delegate":        "委派",
		"add_approver":    "添加审批人",
		"remove_approver": "移除审批人",
		"expire":          "过期",
		"complete":        "完成",
		"cancel":          "取消",
		"restart":         "重启",
		"comment":         "评论",
		"view":            "查看",
		"remind":          "提醒",
		"escalate":        "升级",
		"auto_approve":    "自动批准",
		"auto_reject":     "自动拒绝",
		"resubmit":        "再次提交",
	}
}

// 审批操作类型多语言支持
var approvalOperationTypeLangMap = map[string]map[string]string{
	"zh": {
		"submit":          "提交审批",
		"approve":         "批准",
		"reject":          "拒绝",
		"return":          "退回修改",
		"withdraw":        "撤销",
		"suspend":         "挂起",
		"resume":          "恢复",
		"transfer":        "转交",
		"delegate":        "委派",
		"add_approver":    "添加审批人",
		"remove_approver": "移除审批人",
		"expire":          "过期",
		"complete":        "完成",
		"cancel":          "取消",
		"restart":         "重启",
		"comment":         "评论",
		"view":            "查看",
		"remind":          "提醒",
		"escalate":        "升级",
		"auto_approve":    "自动批准",
		"auto_reject":     "自动拒绝",
		"resubmit":        "再次提交",
	},
	"en": {
		"submit":          "Submit",
		"approve":         "Approve",
		"reject":          "Reject",
		"return":          "Return",
		"withdraw":        "Withdraw",
		"suspend":         "Suspend",
		"resume":          "Resume",
		"transfer":        "Transfer",
		"delegate":        "Delegate",
		"add_approver":    "Add Approver",
		"remove_approver": "Remove Approver",
		"expire":          "Expire",
		"complete":        "Complete",
		"cancel":          "Cancel",
		"restart":         "Restart",
		"comment":         "Comment",
		"view":            "View",
		"remind":          "Remind",
		"escalate":        "Escalate",
		"auto_approve":    "Auto Approve",
		"auto_reject":     "Auto Reject",
		"resubmit":        "Resubmit",
	},
	"ru": {
		"submit":          "отправить",
		"approve":         "одобрить",
		"reject":          "отклонить",
		"return":          "вернуть",
		"withdraw":        "отозвать",
		"suspend":         "приостановить",
		"resume":          "возобновить",
		"transfer":        "передать",
		"delegate":        "делегировать",
		"add_approver":    "добавить утверждающего",
		"remove_approver": "удалить утверждающего",
		"expire":          "истечь",
		"complete":        "завершить",
		"cancel":          "отменить",
		"restart":         "перезапустить",
		"comment":         "комментарий",
		"view":            "просмотр",
		"remind":          "напомнить",
		"escalate":        "эскалация",
		"auto_approve":    "автоодобрение",
		"auto_reject":     "автоотклонение",
		"resubmit":        "повторно отправить",
	},
}

// ApprovalActionLangLabels 返回指定语言的审批操作类型标签
func ApprovalActionLangLabels(langs ...string) (map[string]string, bool) {
	return lang.MatchContent(approvalOperationTypeLangMap, langs)
}

// String 返回字符串
func (a ApprovalAction) String() string {
	return string(a)
}

// Label 返回标签
func (a ApprovalAction) Label() string {
	return ApprovalActionLabels()[a.String()]
}

// LangLabel 返回多语言标签
func (a ApprovalAction) LangLabel(langs ...string) (string, bool) {
	labels, b := ApprovalActionLangLabels(langs...)
	if !b {
		return a.Label(), false
	}
	return labels[a.String()], true
}

// Values 返回所有值
func (a ApprovalAction) Values() []string {
	return ApprovalActionValues()
}

// IsUserOperation 判断是否为用户操作（非系统自动操作）
func (a ApprovalAction) IsUserOperation() bool {
	switch a {
	case ApprovalActionAutoApprove, ApprovalActionAutoReject, ApprovalActionExpire:
		return false
	default:
		return true
	}
}

// IsSystemOperation 判断是否为系统操作
func (a ApprovalAction) IsSystemOperation() bool {
	return !a.IsUserOperation()
}

// IsApprovalOperation 判断是否为审批相关操作
func (a ApprovalAction) IsApprovalOperation() bool {
	switch a {
	case ApprovalActionApprove, ApprovalActionReject, ApprovalActionReturn:
		return true
	default:
		return false
	}
}

// IsFlowControlOperation 判断是否为流程控制操作
func (a ApprovalAction) IsFlowControlOperation() bool {
	switch a {
	case ApprovalActionSubmit, ApprovalActionWithdraw, ApprovalActionSuspend,
		ApprovalActionResume, ApprovalActionCancel, ApprovalActionRestart,
		ApprovalActionComplete:
		return true
	default:
		return false
	}
}

// IsManagementOperation 判断是否为管理操作
func (a ApprovalAction) IsManagementOperation() bool {
	switch a {
	case ApprovalActionTransfer, ApprovalActionDelegate, ApprovalActionAddApprover,
		ApprovalActionRemoveApprover, ApprovalActionEscalate:
		return true
	default:
		return false
	}
}
