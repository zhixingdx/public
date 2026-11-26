package ann

// BizCallAPI 为租户内有唯一业务键的表生成操作。缺省未命中返回 NotFound。
type BizCallAPI struct {
	Manage BizCallActions `json:"manage,omitempty"`
	Portal BizCallActions `json:"portal,omitempty"`
	// InitOnMissing 保留配置表的首次读取初始化和首次编辑创建语义。
	InitOnMissing bool `json:"initOnMissing,omitempty"`
	// StructReply 只返回指定的 Struct 配置正文。留空时复用本端详情投影。
	StructReply string `json:"structReply,omitempty"`
}

// BizCallActions 分端声明读取与写入。编辑白名单为空时不生成编辑接口。
type BizCallActions struct {
	Detail     bool     `json:"detail,omitempty"`
	EditFields []string `json:"editFields,omitempty"`
	HTTP       bool     `json:"http,omitempty"`
}
