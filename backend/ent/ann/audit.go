package ann

// AuditConfig 声明源数据动作、事件版本及排除字段，不要求中央逐动作注册。
type AuditConfig struct {
	// ActionPrefix 留空时使用服务名与实体 snake_case，例如 iam.private_role。
	ActionPrefix string `json:"actionPrefix"`
	// PolicyVersion 表达源事件版本，零值默认 1，不作为中央配置查找键。
	PolicyVersion int32 `json:"policyVersion"`
	// UpdatePolicyVersion 留空时继承 PolicyVersion，正数可单独覆盖修改事件版本。
	UpdatePolicyVersion int32 `json:"updatePolicyVersion,omitempty"`
	// ExcludeFields 同时从普通增改事件和删除归档排除字段，保持既有配置语义。
	ExcludeFields []string `json:"excludeFields,omitempty"`
	// ChangeExcludeFields 仅从普通增改事件排除字段，删除归档仍保留；凭据始终排除。
	ChangeExcludeFields []string `json:"changeExcludeFields,omitempty"`
}

// DeletionMode 表达源表删除生命周期，与普通修改审计的采集配置分离。
type DeletionMode string

// DeletionModeArchive 可独立使用。永久快照成功后才物理删除，事务提交后确认，恢复由来源业务负责。
const DeletionModeArchive DeletionMode = "archive"
