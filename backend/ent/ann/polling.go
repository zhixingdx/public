package ann

// PollingConfig 显式选择一条轮询流及字段映射；未声明的表不生成轮询能力。
// ID 与 Group 使用稳定字符串键；资格字段仅生成强类型 IN 过滤，不接受任意 SQL。
type PollingConfig struct {
	// WorkflowName 可选覆盖生成的 Workflow/Biz 入口名，缺省为 Poll 加表名。
	WorkflowName string   `json:"workflowName,omitempty"`
	JobName      string   `json:"jobName"`
	ID           string   `json:"id"`
	Group        string   `json:"group"`
	Checked      string   `json:"checked"`
	Ready        string   `json:"ready"`
	FilterFields []string `json:"filterFields"`
}
