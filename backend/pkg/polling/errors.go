package polling

// Error 保存跨层传递的结构化原因；服务边界负责 fault 和多语言，不解析错误文本。
type Error struct {
	Reason string
	Field  string
	Limit  string
}

// NewError 创建不依赖服务或语言包的轮询错误。
func NewError(reason string) *Error { return &Error{Reason: reason} }

// Error 返回稳定诊断键；用户可见文案由 CommonLang 映射。
func (e *Error) Error() string { return "polling_" + e.Reason }
