package ann

// ModerationConfig 声明字段需要文本检测或替换，空场景使用默认关键词。
type ModerationConfig struct {
	Scene string `json:"scene,omitempty"`
}
