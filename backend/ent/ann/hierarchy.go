package ann

// Hierarchy 声明受 SassField 隔离的标准接口下级能力；SassOnlyPortal 表仅可开启 Portal。
type Hierarchy struct {
	Resource string  `json:"resource"`
	Manage   TableFn `json:"manage"`
	Portal   TableFn `json:"portal"`
}
