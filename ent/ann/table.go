// Package ann Table
package ann

import (
	"encoding/json"
	"log"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/zhixingdx/public/api/apicommon/v1"
)

const (
	// BatchMaximumSize 是单次批量新增等通用批量请求的全仓硬上限，避免请求体、SQL 参数与事务锁集合无界增长。
	BatchMaximumSize = 1000
	// BatchEditBulkMaximumPatches 是单次异构批量更新的全仓硬上限，避免 SQL 参数、锁集合与请求体无界增长。
	BatchEditBulkMaximumPatches = BatchMaximumSize
)

// BatchEditBulkMode 定义 BatchEdit 在标准 Ent 更新与派生表批量更新之间的选择策略。
type BatchEditBulkMode string

const (
	// BatchEditBulkModeEnt 强制逐条执行 Ent UpdateOne；空值也按该模式处理，以保持既有 Hook 语义。
	BatchEditBulkModeEnt BatchEditBulkMode = "ent"
	// BatchEditBulkModeAuto 在静态能力与运行时 Hook 都兼容时使用批量更新，否则安全降级到 Ent UpdateOne。
	BatchEditBulkModeAuto BatchEditBulkMode = "auto"
	// BatchEditBulkModeRequired 要求必须使用批量更新；检测到不兼容能力时拒绝生成或执行，避免性能静默退化。
	BatchEditBulkModeRequired BatchEditBulkMode = "required"
)

// Effective 将零值归一化为 Ent 模式，避免未配置的历史表被自动切换到绕过 Hook 的路径。
func (m BatchEditBulkMode) Effective() BatchEditBulkMode {
	if m == "" {
		return BatchEditBulkModeEnt
	}
	return m
}

// IsValid 判断模式是否属于生成器支持的稳定契约。
func (m BatchEditBulkMode) IsValid() bool {
	switch m.Effective() {
	case BatchEditBulkModeEnt, BatchEditBulkModeAuto, BatchEditBulkModeRequired:
		return true
	default:
		return false
	}
}

// IsBulkEnabled 判断模式是否允许生成 MySQL/PostgreSQL 派生表批量更新能力。
func (m BatchEditBulkMode) IsBulkEnabled() bool {
	switch m.Effective() {
	case BatchEditBulkModeAuto, BatchEditBulkModeRequired:
		return true
	default:
		return false
	}
}

// BatchAPI 控制批量写能力是否通过 Manage RPC/HTTP 对外暴露；零值默认不暴露。
// 该配置只负责接口暴露，Biz/Data 仍保持生成，避免与 Custom 的实现归属语义混淆。
type BatchAPI struct {
	Add  bool `json:"add"`
	Edit bool `json:"edit"`
}

// Table 数据库表的注解
type Table struct {
	BizCall *BizCallAPI `json:"bizCall,omitempty"`
	// WriteGuard 显式要求标准新增与批量编辑在同一事务中先调用实体写守卫；缺省不改变生成行为。
	// 当前仅支持具有独立主键的单租户实体，不支持 CrossTenantData 和树聚合。
	WriteGuard bool `json:"writeGuard,omitempty"`
	// DeletionMode 显式选择永久归档后物理删除。空值保留原业务删除方式，任何软删字段均不会隐式启用归档。
	DeletionMode DeletionMode `json:"deletionMode,omitempty"`
	// Audit 启用普通增改与硬删元数据采集。只设置 DeletionModeArchive 时仅在删除前归档。
	Audit *AuditConfig `json:"audit,omitempty"`
	// Polling 为单个实体显式启用通用轮询查询与独立检查事实写入。
	Polling *PollingConfig `json:"polling,omitempty"`

	// 表名 用于展示
	TableName string `json:"tableName"`

	// 默认每页数量 默认 20
	DefaultPageSize int `json:"defaultPageSize"`
	// 最大每页数量 默认 1000
	MaxPageSize int `json:"maxPageSize"`
	// RequirePositivePageSize 显式拒绝 pageSize=0；缺省保留既有“不加 LIMIT”兼容语义。
	RequirePositivePageSize bool `json:"requirePositivePageSize,omitempty"`
	// 批量请求最大元素数，0 表示使用全仓硬上限；正数只能进一步收窄，不能超过全仓硬上限。
	MaxBatchSize int `json:"maxBatchSize"`
	// BatchEditBulkMode 控制 BatchEdit 使用 Ent、自动降级或强制批量更新；零值保持 Ent 兼容路径。
	// Auto 会在运行时检测手写 Ent Hook，Required 则在发现 Hook 时拒绝执行，避免静默绕过副作用。
	BatchEditBulkMode BatchEditBulkMode `json:"batchEditBulkMode"`
	// BatchEditBulkMaxSize 限制批量直写的最大补丁数；0 使用全仓默认上限，正数只能进一步收窄。
	// 该字段与通用 Multi/BatchAdd 的 MaxBatchSize 分离，避免两类操作的容量策略互相耦合。
	BatchEditBulkMaxSize int `json:"batchEditBulkMaxSize"`
	// 默认排序
	DefaultSortOrder []*apicommon.SortOrder `json:"defaultSort"`

	// 是否启用返回所有记录 默认不启用
	EnableQueryAll bool `json:"enableQueryAll"`
	// 缓存 缓存以单表在 data 层进行存储，边数据，为防边循环和与data数据一致，只取一层边关系数据
	Cache TableCache `json:"cache"`
	// 聚合与分组
	Aggregates []Aggregate `json:"aggregate"`
	// 树的父级字段名 为空则不启用 值为字段名
	TreeParent string `json:"treeParent"`
	// 边
	WithEdge map[string]WithConfig `json:"withEdge"`
	// json 字段的 ids 带上对应表的数据
	WithTableByJSON map[string]WithTableByJSON `json:"withTableByJSON"`
	// 服务
	WithServices []WithService `json:"withServices"`
	// 自定义方法配置
	Custom TableCustom `json:"custom"`
	// ServiceOnlyCustom 只关闭自动 Service/Proto/权限入口，保留生成的 Biz/Data 供受控入口复用。
	ServiceOnlyCustom TableFn `json:"serviceOnlyCustom"`
	// ManageBatchAPI 显式控制 Manage 批量新增、批量编辑接口；默认关闭以收窄高放大写入口。
	ManageBatchAPI BatchAPI `json:"manageBatchApi"`
	// ManageNoPermission 仅给常规 Manage RPC 的指定动作生成免权限注释，供操作项清单排除既有业务授权入口。
	// 它不改变服务中间件白名单、不开放额外 RPC；缓存及 My/Public/Accessible 等专属接口仍使用各自配置。
	ManageNoPermission TableFn `json:"manageNoPermission"`
	// 是否生成 Http 接口 默认为空不生成 仅提供 rpc 如有值则为接口前缀
	HTTPAPI string `json:"httpApi"`
	// Portal 提供给前端的接口, 通过是部分用于管理的字段不返回给前端
	PortalAPI TableFn `json:"portalApi"`
	// 加密的密钥 Key 如不填则使用默认
	SecretKey string `json:"secretKey"`

	OpProtectConfig OpProtectConfig `json:"operationProtectionConfig"`

	MutationConfig MutationConfig `json:"mutationConfig"`

	// 是否是日志操作表
	IsOpLogTable bool `json:"isOpLogTable"`

	// 权限配置的ID
	PermissionsID string `yaml:"permissionsParentID"`

	// 自定义字段 业务逻辑字段 不存在于数据库 仅用于代码生成
	CustomFields []CustomField `json:"customFields"`

	// 切面编程 用于在代码生成时 对代码进行切面处理
	// 支持在  service biz 层对自动生成的每个 fn 进行处理
	ServiceBefore       TableFn `json:"serviceBefore"`
	ServiceAfter        TableFn `json:"serviceAfter"`
	ServicePortalBefore TableFn `json:"servicePortalBefore"`
	ServicePortalAfter  TableFn `json:"servicePortalAfter"`
	BizBefore           TableFn `json:"bizBefore"`
	BizAfter            TableFn `json:"bizAfter"`
	// DataScopeWhere 在生成的查询与操作条件出口调用实体自定义 Data 范围谓词切面。
	// 它适合租户、Realm 等必须覆盖读取与按条件写入的底层数据边界。
	DataScopeWhere bool `json:"dataScopeWhere"`
	// Hierarchy 开启受 SassField 隔离入口的下级 CRUD；SassOnlyPortal 表仅可开启 Portal，零值关闭。
	Hierarchy *Hierarchy `json:"hierarchy,omitempty"`

	// 是否支持对外提供跨租户查询的方法
	CrossTenantFn       CrossTenantFn `json:"crossTenantFn"`
	PortalCrossTenantFn CrossTenantFn `json:"portalCrossTenantFn"`
	// 自定义实现跨租户 biz 层逻辑
	CrossTenantBizCustomFn CrossTenantFn `json:"crossTenantBizCustomFn"`
	// data 层是否提供跨租户操作( 跨租户操作添加/编辑/删除/停启用相当敏感,只通过这个开关，在 data 层自动生成对应的方法,再手动在 biz 层进行调用处理)
	CrossTenantData bool `json:"crossTenantData"`

	// 是否强制生成 add proto 请求参数
	ForceAddProtoRequest bool `json:"forceAddProtoRequest"`
	// 是否强制生成 edit proto 请求参数
	ForceEditProtoRequest bool `json:"forceEditProtoRequest"`
	// Portal 层是否强制生成 add proto 请求参数
	ForceAddPortalProtoRequest bool `json:"forceAddPortalProtoRequest"`
	// Portal 层是否强制生成 edit proto 请求参数
	ForceEditPortalProtoRequest bool `json:"forceEditPortalProtoRequest"`

	// 是否在 data 层生成 with tx 方法 (目前为 add 和 edit 方法)
	WithTx bool `json:"dbWithTx"`

	// 草稿发布配置
	// 用于配置草稿表到正式表的自动同步规则
	Draft *DraftConfig `json:"draft,omitempty"`
	// VersionedResource 配置稳定 Head 与 Version 内容事实源的聚合生成能力。
	// 它与 DraftConfig 的复制模型语义相反，因此生成器会拒绝同时配置。
	VersionedResource *VersionedResourceConfig `json:"versionedResource,omitempty"`
	// ResourceMetadata 配置 owner 实体对专属 metadata 表的显式生成能力。
	// OwnerField/ScopeField 只由 metadata Schema 的 ResourceMetadataStorage 声明，避免双份配置漂移。
	ResourceMetadata *ResourceMetadataConfig `json:"resourceMetadata,omitempty"`

	PublicReadableFn TableReadFn `json:"publicReadableFn"`
	// MyOwnerFn 配置 Manage /my/ 作用域；一张表只能配置一个 MyOwner 字段。
	// Add/Edit 需同时声明 MyOwnerWrite，并由领域层实现同名 MyAdd/MyEdit 用例。
	// Del 显式生成整批个人删除，要求原生字符串主键与不可变的原生字符串 MyOwner、租户归属；每次最多 100 个 ID。
	MyOwnerFn TableFn `json:"myOwnerFn"`
	// MyOwnerWrite 声明 Add/Edit 入口允许填写的字段；未声明的管理字段不会随实体扩展自动开放。
	MyOwnerWrite *MyOwnerWriteConfig `json:"myOwnerWrite,omitempty"`
	// AccessibleFn 标识 /accessible/ 作用域下的租户可见范围只读接口
	// 当前仅支持具备 visibilityScope 与 visibleTenants 字段的表
	AccessibleFn TableReadFn `json:"accessibleFn"`
	// AccessibleServiceCustom 标识由 Plus 自定义承接的 Accessible Service 方法。
	// Data/Biz 仍按 AccessibleFn 自动生成，只有 Proto/Service/Server 跳过对应方法，供安全投影等非标准响应复用底层可见性查询。
	AccessibleServiceCustom TableReadFn `json:"accessibleServiceCustom"`

	// 获取 OwnerID 的方法。默认: ownerID, ownerErr := meta.GetStaffID(ctx)；生成代码随后会判断 ownerErr 并 HandleMetadataError。
	// 自定义片段须与默认一致使用 ownerErr（勿用 err），否则与 usecase 返回的 error 复用会触发类型推断错误。
	OwnerIDFn string `json:"ownerIDFn"`

	PublicReadablePortalFn TableReadFn `json:"publicReadablePortalFn"`
	// MyOwnerFn 标识 /my/ 作用域下的归属判断字段（当前用户与数据所有者匹配）
	// 一张表只能配置一个 MyOwner 字段
	MyOwnerPortalFn TableFn `json:"myOwnerPortalFn"`

	// 获取 OwnerID 的方法。默认: ownerID, ownerErr := meta.GetStaffID(ctx)；生成代码随后会判断 ownerErr 并 HandleMetadataError。
	// 自定义片段须与默认一致使用 ownerErr（勿用 err），否则与 usecase 返回的 error 复用会触发类型推断错误。
	OwnerIDPortalFn string `json:"ownerIDPortalFn"`

	// ManageReadWithoutMask 配置管理侧「不脱敏」只读 RPC（与常规读接口同请求/响应消息类型，phone/mail 等为明文）。
	// RPC 由 yimoka 写入 auto_service.proto（service {Repo}Manage），实现见 internal/service/auto_read_without_mask.go。
	// 仅生成 gRPC、不配置 google.api.http；后续须在此类 RPC 上强制接入审计。
	// 仅当表中存在 ann.Field.MaskEncrypt（等敏感掩码）字段时允许配置；与 Manage 层对应读接口同源，若某读能力为「用户自定义」则不同名生成该 WithoutMask 方法。
	ManageReadWithoutMask TableReadFn `json:"manageReadWithoutMask"`

	// SearchConfig 搜索配置
	// 配置后自动生成 OpenSearch 同步和搜索相关代码
	SearchConfig *SearchConfig `json:"searchConfig,omitempty"`
}

// SearchConfig 表级搜索配置
// 用于配置 OpenSearch 索引映射、数据同步和搜索接口生成
type SearchConfig struct {
	// IndexName OpenSearch 索引名（多表共用同一索引时需指定相同名称）
	// 默认为空，使用表名的 snake_case 形式
	IndexName string `json:"indexName,omitempty"`

	// TypeField 类型标识字段名（多表共用索引时，区分数据来源）
	// 值自动填充为表名的 snake_case
	TypeField string `json:"typeField,omitempty"`

	// ExcludeFields 排除字段列表（默认同步全部字段，此处指定不同步的字段）
	ExcludeFields []string `json:"excludeFields,omitempty"`

	// ExtraFields 额外字段（不在数据库表中，但需要同步到索引的字段）
	ExtraFields []SearchExtraField `json:"extraFields,omitempty"`

	// KeywordFields 全文搜索字段配置（支持权重）
	// key: 字段名，value: 权重（正整数，越大权重越高，默认 1）
	KeywordFields map[string]int `json:"keywordFields,omitempty"`

	// ManageSearch 是否生成管理端搜索接口
	ManageSearch bool `json:"manageSearch,omitempty"`

	// BizBefore 是否在搜索 UseCase 调用 Repo 前执行 search{Entity}Before 业务钩子。
	BizBefore bool `json:"bizBefore,omitempty"`

	// PortalSearch 是否生成门户端搜索接口
	PortalSearch bool `json:"portalSearch,omitempty"`

	// SyncHook 是否启用 Hook 异步同步（默认 true）
	// 设为 false 时仅生成同步接口，不自动 Hook 同步
	SyncHook *bool `json:"syncHook,omitempty"`

	// ManageSyncAPI 是否在 Manage 层生成同步 RPC（默认 false）
	// 租户可调用的同步接口，谨慎开启，避免影响性能
	ManageSyncAPI bool `json:"manageSyncAPI,omitempty"`

	// WorkflowSync 是否在 Workflow 层生成同步 Activity（默认 true）
	// 运维通过 Temporal 触发，不限租户，全量同步所有数据
	WorkflowSync *bool `json:"workflowSync,omitempty"`

	// GenerateToEntityFn 是否生成搜索结果转表实体的转换函数
	GenerateToEntityFn bool `json:"generateToEntityFn,omitempty"`
}

// IsWorkflowSyncEnabled 是否启用 Workflow 同步（默认 true）
func (c *SearchConfig) IsWorkflowSyncEnabled() bool {
	if c.WorkflowSync == nil {
		return true
	}
	return *c.WorkflowSync
}

// IsSyncHookEnabled 是否启用 Hook 同步（默认 true）
func (c *SearchConfig) IsSyncHookEnabled() bool {
	if c.SyncHook == nil {
		return true
	}
	return *c.SyncHook
}

// GetIndexName 获取索引名（默认为 snake_case 表名）
func (c *SearchConfig) GetIndexName(tableName string) string {
	if c == nil {
		return toSnakeCase(tableName)
	}
	indexName := strings.TrimSpace(c.IndexName)
	if indexName != "" {
		return indexName
	}
	if c.IndexName != "" {
		return ""
	}
	return toSnakeCase(tableName)
}

// IsExcluded 检查字段是否被排除
func (c *SearchConfig) IsExcluded(fieldName string) bool {
	if c == nil {
		return false
	}
	fieldName = strings.TrimSpace(fieldName)
	if fieldName == "" {
		return false
	}
	for _, f := range c.ExcludeFields {
		if strings.TrimSpace(f) == fieldName {
			return true
		}
	}
	return false
}

// SearchExtraField 搜索索引的额外字段
type SearchExtraField struct {
	// Name 字段名
	Name string `json:"name"`
	// Type OpenSearch 字段类型（keyword/text/integer/long/boolean/date 等）
	Type string `json:"type"`
	// Properties OpenSearch 对象字段的子字段类型，如 {"view":"long","favorite":"long"}
	Properties map[string]string `json:"properties,omitempty"`
	// Value 固定值（如 "article"）
	Value string `json:"value"`
	// PbIndex protobuf 字段序号，遵循与表字段相同的编号规范
	PbIndex int `json:"pbIndex"`
	// PbType 自定义 protobuf 类型；为空时由 Type 推断
	PbType string `json:"pbType,omitempty"`
	// GoType 自定义 Go 类型；为空时由 Type 推断
	GoType string `json:"goType,omitempty"`
	// Comment 字段注释
	Comment string `json:"comment,omitempty"`
}

func (f SearchExtraField) IsObjectMapping() bool {
	return len(f.Properties) > 0
}

func (f SearchExtraField) IsFilterable() bool {
	return strings.TrimSpace(f.Type) != "" && !f.IsObjectMapping()
}

func (f SearchExtraField) HasFixedValue() bool {
	return strings.TrimSpace(f.Value) != ""
}

// toSnakeCase 将 PascalCase 转换为 snake_case
func toSnakeCase(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(c-'A'+'a'))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}

// DraftConfig 草稿发布配置
// 遵循"约定优于配置"原则：
// - 同名字段自动映射
// - AutoCreate/AutoUpdate 标记的字段自动排除
// - 只需配置例外情况
type DraftConfig struct {
	// TargetTable 目标正式表（Ent Schema 名称）
	// 例如: "Course" 表示草稿发布到 Course 表
	TargetTable string `json:"targetTable"`

	// ExcludeFields 排除的字段列表（不同步到正式表）
	// 草稿特有字段如 status, sourceID 等
	// 注意: AutoCreate/AutoUpdate 的字段会自动排除，无需配置
	ExcludeFields []string `json:"excludeFields,omitempty"`

	// FieldMapping 字段映射（源字段 -> 目标字段）
	// 用于字段名不同的情况
	// 例如: {"staffID": "publisherID"} 表示草稿的 staffID 映射到正式表的 publisherID
	FieldMapping map[string]string `json:"fieldMapping,omitempty"`
}

// MyOwnerWriteConfig 声明个人创建、编辑的字段允许集，领域层继续实现同名 MyAdd/MyEdit 用例。
type MyOwnerWriteConfig struct {
	// AddFields 使用 Ent 字段名，只允许既有 Add 请求中的普通字段。
	AddFields []string `json:"addFields,omitempty"`
	// EditFields 使用 Ent 字段名；操作主键与允许集合字段的更新辅助键由生成器补齐。
	EditFields []string `json:"editFields,omitempty"`
}

// VersionedResourceConfig 声明稳定资源 Head 与不可变历史 Version 的生成合同。
type VersionedResourceConfig struct {
	// VersionTable 指定与当前 Head 配对的 Ent Version Schema 名称。
	VersionTable string `json:"versionTable"`
	// DraftReplyPbIndex 显式指定 Head 草稿投影编号；零值沿用既有 20。
	DraftReplyPbIndex int `json:"draftReplyPbIndex,omitempty"`
	// PublishedReplyPbIndex 显式指定 Head 发布投影编号；零值沿用既有 21。
	PublishedReplyPbIndex int `json:"publishedReplyPbIndex,omitempty"`
	// ReservedHeadReplyPbIndexes 保留已移出 Head 响应的旧字段号，避免被新字段复用。
	ReservedHeadReplyPbIndexes []int `json:"reservedHeadReplyPbIndexes,omitempty"`
	// PayloadFields 指定只保存在 Version 中并参与 contentHash 的完整 payload 字段。
	PayloadFields []string `json:"payloadFields"`
	// PayloadOmitEmptyFields 为新增 JSON 对象字段保留缺省值的历史摘要。
	PayloadOmitEmptyFields []string `json:"payloadOmitEmptyFields,omitempty"`
	// StableFields 指定保留在 Head 上且创建后不可变的业务稳定标识。
	StableFields []string `json:"stableFields,omitempty"`
	// MutableMetadataFields 指定 Version 上不参与 contentHash 的可变控制元数据。
	MutableMetadataFields []string `json:"mutableMetadataFields,omitempty"`
	// ActivatePublished 允许管理员将当前发布指针切回已发布的历史版本。
	ActivatePublished bool `json:"activatePublished,omitempty"`
	// PublishBefore 声明 Biz 层必须在短事务前执行纯本地发布预检钩子。
	PublishBefore bool `json:"publishBefore,omitempty"`
}

// IsExcluded 检查字段是否被排除
func (c *DraftConfig) IsExcluded(fieldName string) bool {
	if c == nil {
		return false
	}
	for _, f := range c.ExcludeFields {
		if f == fieldName {
			return true
		}
	}
	return false
}

// GetTargetFieldName 获取目标字段名（处理字段映射）
func (c *DraftConfig) GetTargetFieldName(sourceFieldName string) string {
	if c == nil || c.FieldMapping == nil {
		return sourceFieldName
	}
	if targetName, ok := c.FieldMapping[sourceFieldName]; ok {
		return targetName
	}
	return sourceFieldName
}

// HasFieldMapping 检查是否有字段映射
func (c *DraftConfig) HasFieldMapping(sourceFieldName string) bool {
	if c == nil || c.FieldMapping == nil {
		return false
	}
	_, ok := c.FieldMapping[sourceFieldName]
	return ok
}

// Place 代码生成的层
type Place string

const (
	// PlaceProto proto 层
	PlaceProto Place = "proto"
	// PlaceService service 层
	PlaceService Place = "service"
	// PlaceBiz biz 层
	PlaceBiz Place = "biz"
	// PlaceData data 层
	PlaceData Place = "data"
)

const tableNameKey = "Table"

// Name _
func (Table) Name() string {
	return tableNameKey
}

// GetTableConfig 获取表配置
func GetTableConfig(node *gen.Type) *Table {
	ann := node.Annotations[tableNameKey]
	if ann == nil {
		return &Table{}
	}
	data, err := json.Marshal(ann)
	if err != nil {
		log.Fatal(err)
	}
	var conf Table
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}
	if conf.DefaultPageSize == 0 {
		conf.DefaultPageSize = 20
	}
	if conf.MaxPageSize == 0 {
		conf.MaxPageSize = 1000
	}
	return &conf
}

// TableCache 表缓存配置
type TableCache struct {
	// 详情缓存 支持配置多个字段 请确保字段的值唯一
	Detail []string `json:"detail"`

	DetailBizBefore map[string]TableCacheFn `json:"detailBizBefore"`
	DetailBizAfter  map[string]TableCacheFn `json:"detailBizAfter"`

	DetailBefore map[string]TableCacheFn `json:"detailBefore"`
	DetailAfter  map[string]TableCacheFn `json:"detailAfter"`
	// 详情缓存免鉴权配置（默认仅放开 Get）
	DetailNoAuth   map[string]TableCacheFn `json:"detailNoAuth"`
	PortalDetailFn map[string]TableCacheFn `json:"portalDetailFn"`
	// Portal 详情缓存免鉴权配置（默认仅放开 Get）
	PortalDetailNoAuth map[string]TableCacheFn `json:"portalDetailNoAuth"`
	// Portal 详情缓存免权限配置（仍需登录，默认仅放开 Get）
	PortalDetailNoPermission map[string]TableCacheFn `json:"portalDetailNoPermission"`
	PortalDetailBefore       map[string]TableCacheFn `json:"portalDetailBefore"`
	PortalDetailAfter        map[string]TableCacheFn `json:"portalDetailAfter"`

	// 启用查询所有缓存
	QueryAll bool `json:"queryAll"`
	// 树缓存接口免鉴权配置（Manage）
	TreeCacheNoAuth bool `json:"treeCacheNoAuth"`
	// 树缓存接口免鉴权配置（Portal）
	PortalTreeCacheNoAuth bool `json:"portalTreeCacheNoAuth"`
	// 查询所有缓存免鉴权配置（默认仅放开 Get）
	QueryAllNoAuth    TableCacheFn `json:"queryAllNoAuth"`
	QueryAllBizBefore TableCacheFn `json:"queryAllBizBefore"`
	QueryAllBizAfter  TableCacheFn `json:"queryAllBizAfter"`

	QueryAllBefore TableCacheFn `json:"queryAllBefore"`
	QueryAllAfter  TableCacheFn `json:"queryAllAfter"`
	// Portal层查询所有缓存 启用的方法 all 只有三个方法 get set del
	PortalQueryAllFn TableCacheFn `json:"portalQueryAll"`
	// Portal 查询所有缓存免鉴权配置（默认仅放开 Get）
	PortalQueryAllNoAuth TableCacheFn `json:"portalQueryAllNoAuth"`
	PortalQueryAllBefore TableCacheFn `json:"portalQueryAllBefore"`
	PortalQueryAllAfter  TableCacheFn `json:"portalQueryAllAfter"`

	// 查询缓存 支持配置多个字段
	Query          []string                `json:"query"`
	QueryBizBefore map[string]TableCacheFn `json:"queryBizBefore"`
	QueryBizAfter  map[string]TableCacheFn `json:"queryBizAfter"`

	QueryBefore map[string]TableCacheFn `json:"queryBefore"`
	QueryAfter  map[string]TableCacheFn `json:"queryAfter"`
	// 查询缓存免鉴权配置（默认仅放开 Get）
	QueryNoAuth map[string]TableCacheFn `json:"queryNoAuth"`
	// Portal层查询缓存 启用的方法
	PortalQueryFn map[string]TableCacheFn `json:"portalQueryFn"`
	// Portal 查询缓存免鉴权配置（默认仅放开 Get）
	PortalQueryNoAuth map[string]TableCacheFn `json:"portalQueryNoAuth"`
	PortalQueryBefore map[string]TableCacheFn `json:"portalQueryBefore"`
	PortalQueryAfter  map[string]TableCacheFn `json:"portalQueryAfter"`

	// 是否不开启穿透保护
	// 穿透保护查询不到的数据 会缓存一个空值
	// 详情缓存 返回 404 列表缓存 返回 空数组
	NotPreventPenetration bool `json:"notPreventPenetration"`
}

// TableCacheFn _
type TableCacheFn struct {
	Get  bool
	MGet bool
	Set  bool
	MSet bool
	Del  bool
}

// WithConfig _
type WithConfig struct {
	Detail bool
	Rows   bool
	Add    bool
	Edit   bool
	Query  bool
}

// WithTableByJSON json 字段的 ids 带上对应表的数据
type WithTableByJSON struct {
	Table  string
	Detail bool
	Rows   bool
}

// TableCustom _
type TableCustom struct {
	Service TableFn
	Biz     TableFn
	Data    TableFn
}

// TableFn _
type TableFn struct {
	Add       bool `json:"add"`
	BatchAdd  bool `json:"batchAdd"`
	Edit      bool `json:"edit"`
	BatchEdit bool `json:"batchEdit"`
	// Current 标识 VersionedResource Head 的 Portal 当前发布版本读取接口。
	Current  bool `json:"current"`
	Detail   bool `json:"detail"`
	Multi    bool `json:"multi"`
	QueryOne bool `json:"queryOne"`
	Query    bool `json:"query"`
	Count    bool `json:"count"`
	List     bool `json:"list"`
	// 以下开关仅供 Portal 生成需要独立权限和必达审计的明文读取，普通读取不会隐式开放。
	DetailWithoutMask   bool `json:"detailWithoutMask"`
	MultiWithoutMask    bool `json:"multiWithoutMask"`
	QueryOneWithoutMask bool `json:"queryOneWithoutMask"`
	QueryWithoutMask    bool `json:"queryWithoutMask"`
	ListWithoutMask     bool `json:"listWithoutMask"`
	All                 bool `json:"all"`
	DelOne              bool `json:"delOne"`
	Del                 bool `json:"del"`
	EnableOne           bool `json:"enableOne"`
	DisableOne          bool `json:"disableOne"`
	Enable              bool `json:"enable"`
	Disable             bool `json:"disable"`
	Tree                bool `json:"tree"`
	// MoveTree 标识树节点移动并排序方法。PortalAPI 中默认关闭，需要门户端拖拽时显式打开。
	MoveTree bool `json:"moveTree"`
}
type TableReadFn struct {
	Detail   bool `json:"detail"`
	Multi    bool `json:"multi"`
	QueryOne bool `json:"queryOne"`
	Query    bool `json:"query"`
	Count    bool `json:"count"`
	List     bool `json:"list"`
	All      bool `json:"all"`
	Tree     bool `json:"tree"`
}

// WithServiceEdgeCondition _
type WithServiceEdgeCondition struct {
	Field string
	Value string
}

// WithServiceEdge 带上其他微服务数据
type WithServiceEdge struct {
	Type         string // O2M M2O M2M 暂时支持 M2O
	EdgeName     string // 边名
	Field        string // 字段名
	DataField    string // 增加存放数据的字段名
	Reply        string // 响应字段名
	QueryType    string // 查询类型 缓存 or db 默认缓存
	IsWithDetail bool   // 是否带上详情
	IsWithRows   bool   // 是否带上列表
	Condition    WithServiceEdgeCondition
}

// WithService 带上其他微服务数据
type WithService struct {
	Name      string
	ProtoPath string
	Package   string
	Usecase   string
	Edges     []WithServiceEdge
}

// Aggregate 聚合与分组
type Aggregate struct {
	// 聚合函数名，为空则 GroupBy 组合
	FuName  string
	GroupBy []string
	// 如果 Aggregate 为空，则默认为 count
	Aggregate []AggregateConf
	// 当没有查询条件时 是否缓存数据
	IsAllCache bool
}

// AggregateConf _
type AggregateConf struct {
	Type  AggregateConfType
	Field string
	As    string
}

// AggregateConfType 聚合类型
type AggregateConfType string

const (
	// AggregateTypeCount 聚合类型 count
	AggregateTypeCount AggregateConfType = "count"
	// AggregateTypeSum 聚合类型 sum
	AggregateTypeSum AggregateConfType = "sum"
	// AggregateTypeMean 聚合类型 mean
	AggregateTypeMean AggregateConfType = "mean"
	// AggregateTypeMax 聚合类型 max
	AggregateTypeMax AggregateConfType = "max"
	// AggregateTypeMin 聚合类型 min
	AggregateTypeMin AggregateConfType = "min"
)

// OpProtectConfig 操作保护的配置 添加 / 编辑 / 删除 / 启用 / 禁用 根据此配置在 biz 层进行插入函数，函数给于用户自行编写
type OpProtectConfig struct {
	// 添加配置
	Add bool `json:"add"`
	// 编辑配置
	Edit bool `json:"edit"`
	// 删除配置
	Del bool `json:"del"`
	// 启用配置
	Enable bool `json:"enable"`
	// 禁用配置
	Disable bool `json:"disable"`
	// Portal 启用配置
	PortalEnable bool `json:"portalEnable"`
	// Portal 禁用配置
	PortalDisable bool `json:"portalDisable"`
}

// MutationConfig 突变配置
type MutationConfig struct {
	// 启用
	Enable bool `json:"enable"`
	// 操作记录表 不填则不启用 需手动创建表,会对表进行需要的字段进行检查
	OpLogTable string `json:"opLogTable"`
	// 操作者 获取语句 必须以 operator 开头 例如 operator := meta.GetUserID(ctx)
	OperatorCode string `json:"operatorCode"`
	// TODO 同步到搜索

}

// CustomField 自定义字段配置
// 用于定义业务逻辑字段，这些字段不存在于数据库中，仅用于代码生成
type CustomField struct {
	// Name 字段名称，用于生成结构体字段名和 JSON 标签
	Name string `json:"name"`

	// Comment 字段注释，用于生成代码注释和 proto 文件中的注释
	Comment string `json:"comment"`

	// PbImport protobuf 导入路径，用于生成 proto 文件中的导入路径
	PbImport string `json:"pbImport"`

	// PbType protobuf 类型定义，用于生成 proto 文件中的字段类型
	PbType string `json:"pbType"`

	// PbIndex protobuf 字段序号，用于生成 proto 文件中的字段编号
	PbIndex int `json:"pbIndex"`

	// GoType Go 语言类型定义，用于生成 Go 结构体中的字段类型
	GoType string `json:"goType"`

	// GoTypeToPb Go 类型转 protobuf 类型的处理代码，如无则直接使用 = 赋值
	// 例如：当需要特殊转换逻辑时，可以指定转换代码
	GoTypeToPb string `json:"goTypeToPb"`

	// PbTypeToGo protobuf 类型转 Go 类型的处理代码，如无则直接使用 = 赋值
	// 例如：当需要特殊转换逻辑时，可以指定转换代码
	PbTypeToGo string `json:"pbTypeToGo"`

	// 是否在 Add 接口的请求中包含此字段
	AddRequest bool `json:"addRequest"`

	// 是否在 Edit 接口的请求中包含此字段
	EditRequest bool `json:"editRequest"`

	// DetailReply 是否在详情接口的响应中包含此字段
	DetailReply bool `json:"detailReply"`

	// QueryReply 是否在查询接口的响应中包含此字段
	QueryReply bool `json:"queryReply"`

	// PortalAddRequest 是否在 Portal 层 Add 接口的请求中包含此字段
	PortalAddRequest bool `json:"portalAddRequest"`

	// PortalEditRequest 是否在 Portal 层 Edit 接口的请求中包含此字段
	PortalEditRequest bool `json:"portalEditRequest"`

	// PortalDetailReply 是否在 Portal 层详情接口的响应中包含此字段
	PortalDetailReply bool `json:"portalDetailReply"`

	// PortalQueryReply 是否在 Portal 层查询接口的响应中包含此字段
	PortalQueryReply bool `json:"portalQueryReply"`

	// Biz 是否在业务层（biz）生成该字段，true 表示在 biz 层结构体中包含此字段
	Biz bool `json:"biz"`
	// BizDetail 是否在业务层详情接口获取值
	BizDetail bool `json:"bizDetail"`
	// BizQuery 是否在业务层查询接口获取值
	BizQuery bool `json:"bizQuery"`
	// BizGetCode 业务层查询接口获取值的代码
	BizGetCode string `json:"bizGetCode"`
	// DataGetCode 字段在数据层获取值的代码
	// 用于在数据层生成该字段的值，通常用于计算字段或从其他数据源获取数据
	DataGetCode string `json:"dataGetCode"`
}

// 是否支持对外提供跨租户查询的方法
type CrossTenantFn struct {
	Detail   bool `json:"detail"`
	Multi    bool `json:"multi"`
	QueryOne bool `json:"queryOne"`
	Query    bool `json:"query"`
	Count    bool `json:"count"`
	List     bool `json:"list"`
}

// HasAny 检查是否有任何一个字段为 true
func (f CrossTenantFn) HasAny() bool {
	return f.Detail || f.Multi || f.QueryOne || f.Query || f.Count || f.List
}
