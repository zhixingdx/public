// Package ann 提供 Yimoka 使用的 Ent Schema 注解合同。
package ann

const (
	// ResourceMetadataStorageNameKey 是 metadata Schema 存储归属注解的稳定名称。
	ResourceMetadataStorageNameKey = "ResourceMetadataStorage"
)

// ResourceMetadataConfig 声明 owner 实体对专属 metadata 表的显式生成配置。
// 存储归属字段不在这里重复配置，由 metadata Schema 的 ResourceMetadataStorage 唯一提供。
type ResourceMetadataConfig struct {
	// Table 指向同一 Ent graph 中的 metadata Schema 名称。
	Table string `json:"table"`
	// RootField 可选声明 owner 的聚合根字段，非根实体不得读写 metadata。
	RootField string `json:"rootField,omitempty"`
	// VersionField 可选声明 owner 的乐观锁版本字段。
	VersionField string `json:"versionField,omitempty"`
	// ManageRPC 显式允许生成内部 Manage gRPC；默认关闭。
	ManageRPC bool `json:"manageRpc,omitempty"`
	// ManageRPCReadOnly 只生成内部读取 RPC，写入由聚合命令统一承接；不能与 ManageRPC 同时开启。
	ManageRPCReadOnly bool `json:"manageRpcReadOnly,omitempty"`
	// Limits 为 owner 的空间、行数和候选 ID 提供有界预算。
	Limits ResourceMetadataLimits `json:"limits"`
}

// ResourceMetadataLimits 冻结实体 metadata 的在线数量预算，零值由生成器替换为安全默认值。
type ResourceMetadataLimits struct {
	// MaxSpacesPerOwner 限制单 owner 可持有的 space 数量。
	MaxSpacesPerOwner int `json:"maxSpacesPerOwner,omitempty"`
	// MaxEntriesPerSpace 限制一次 space 替换后的最大行数。
	MaxEntriesPerSpace int `json:"maxEntriesPerSpace,omitempty"`
	// MaxEntriesPerOwner 限制单 owner 全部 space 的最大行数。
	MaxEntriesPerOwner int `json:"maxEntriesPerOwner,omitempty"`
	// MaxCandidateIDs 限制 metadata 查询向 owner 查询传递的候选 ID 数量。
	MaxCandidateIDs int `json:"maxCandidateIDs,omitempty"`
}

// ResourceMetadataStorage 是 metadata Schema 上 owner/scope 归属的唯一事实源。
type ResourceMetadataStorage struct {
	// OwnerField 是 metadata 行中指向稳定 owner ID 的字段名。
	OwnerField string `json:"ownerField"`
	// ScopeField 是可选隔离字段名；owner 没有 SassField 时必须留空。
	ScopeField string `json:"scopeField,omitempty"`
}

// Name 返回 Ent graph 中 ResourceMetadataStorage 的稳定注解键。
func (ResourceMetadataStorage) Name() string {
	return ResourceMetadataStorageNameKey
}
