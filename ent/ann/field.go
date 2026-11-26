// Package ann Field 字段配置
package ann

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/zhixingdx/public/pkg/utils"
	"github.com/samber/lo"

	"entgo.io/ent/entc/gen"
)

const (
	// maxProtoFieldNumber 是 Protobuf 字段号允许的最大值。
	maxProtoFieldNumber = 1<<29 - 1
	// Protobuf 将 19000-19999 保留给内部实现，业务字段不能占用。
	protoReservedFieldNumberStart = 19000
	protoReservedFieldNumberEnd   = 19999
)

// Field _
type Field struct {
	// PB 的序号
	PbIndex int

	// 表示单字段唯一索引（可忽略 sass 字段，减少 key 长度，并支持某些场景与 sass 字段无关的查询，例如通过域名/appid 反查配置）
	OnlyUnique bool

	// 当传入 zero 值如 “”或 0 时，清除数据库的值为 null (这在唯一索引时是必须的) 处理 字符串 数字 和 时间
	SetNullOnZero bool

	// 特殊字段定义 之所以放在字段里定义是为了可以直接使用字段的 mixin
	// 是否是操作字段 一张表只能有一个操作字段 当 false 默认为 ID
	OperationField bool
	// 是否是 sass 字段 一张表只能有一个 sass 字段
	SassField bool
	// sass 只在 portal 服务使用
	SassOnlyPortal bool
	// 软删除字段 一张表只能有一个软删除字段 并且必须是 bool 类型
	SoftDeleteField bool
	// 行数据开关 停/启用的字段 一张表只能有一个行数据开关字段 并且必须是 bool 类型
	SwitchField bool
	// Portal 行数据开关的字段 只在 Portal 生成开关相关的方法
	SwitchPortalField bool
	// 树根字段 一张树表只能有一个。该字段由后端维护，不暴露在 Add/Edit 请求中。
	TreeRootID bool
	// 树同级排序字段 一张树表只能有一个。用于生成移动节点并重排兄弟节点的 data 层方法。
	TreeSort bool

	//  用于生成 pb 的 validate 文档 https://github.com/bufbuild/protovalidate/blob/main/proto/protovalidate/buf/validate/validate.proto
	// 默认在 Add 和 Edit 时生效
	Validate string
	// 不同接口的 validate 配置
	ValidateMap map[string]string

	// 在查询的 Reply 不返回
	NotQueryReply bool
	// 在详情的 Reply 不返回
	NotDetailReply bool
	// 表示在 Portal 层 Add 时 不需要该字段，例如后台的 备注 字段
	NotPortalAdd bool
	// 表示在 Portal 层 Edit 时 不需要该字段，例如后台的 备注 字段
	NotPortalEdit bool
	// 表示在 Portal 层 Query 时 不需要该字段，例如后台的 备注 字段
	NotPortalQuery bool
	// 在 Portal 层的列表 Reply 不返回
	NotPortalQueryReply bool
	// 在 Portal 层的详情 Reply 不返回
	NotPortalDetailReply bool

	// 查询的配置
	Query FieldQuery
	// 仅是数据层的字段 如辅助加解密的字段
	OnlyData bool
	// 是否只从 meta 中获取, 如 sass 中的 appid
	OnlyMeta bool

	// 自动生成的值，不接受传参 如创建时间
	AutoCreate bool
	// 自动更新的值，不接受传参 如更新时间
	AutoUpdate bool

	// 默认值 用于 生成默认值的方法 ent 默认值不带 ctx，无法获取用户相关的数据
	DefaultFn *FieldFunc
	// 更新值 用于 生成更新值的方法 ent 默认值不带 ctx，无法获取用户相关的数据
	UpdateFn *FieldFunc

	// 当有值时 pb 使用 int64 并根据值是 秒还是毫秒在 service 层转换 time.Time
	PBTimeToType PBTimeType

	// 表示该字段存储需要加密, 仅在 data 层使用
	Encrypt bool
	// 掩码存储, 用于存储敏感数据, 仅在 data 层使用 如手机号
	// 请确保有 字段名 + Cipher 的字段用于存储加密后的数据 并确保 onlyData 为 true
	MaskEncrypt utils.MaskType
	// 独立加密, 用于存储敏感数据, 仅在 data 层使用 如 密码通过生成每一行独立的 nonce 保证安全
	// 请确保有 字段名 + Nonce 的字段用于存储随机生成的 nonce 并确保 onlyData 为 true
	// 请确保有 字段名 + Cipher 的字段用于存储加密后的数据 并确保 onlyData 为 true
	RowIrreversibleEncrypt bool

	// 是否需要 xss 过滤
	XSSFilter bool

	// Moderation 在仓储写入前准备文本，处理模式由业务注入的端口决定。
	Moderation *ModerationConfig

	// 大字段对象存储配置
	// 当字段内容超过阈值时，自动上传到对象存储，数据库只存 URL 引用
	// 仅在 Data 层处理，上层无感知
	ObjectStore *ObjectStoreConfig

	// 多语言字段，为空则不支持多语言，不为空则表示该字段为值字段的多语言值
	// 不为空则必须有对应的默认内容字段，必有为 json 类型 且值与值字段相同
	// 例如：值为 title 时，则必须有 title 字段
	// 本身字段名默认为 xxxI18n
	// 值为 json 类型，key 为语言代码，value 为对应的值
	// 例如：{"zh-CN":"标题","en-US":"title"}
	I18NFor string
	// Portal 层仅本地语言,为 true 在 Portal 层的无 I18N 输入，输出默认内容与和本地语言的内容
	PortalOnlyLocalLang bool
	// PortalI18NEditable 标识 Portal 可读取并编辑完整多语言字段。
	// 仅与 PortalOnlyLocalLang 同时启用时生效，普通列表仍只返回本地语言投影。
	PortalI18NEditable bool
	// PortalI18NEditablePbIndex 指定 Portal 详情中完整多语言字段的独立 Proto 编号。
	// 原 PbIndex 已由 xxxLocal 占用，因此开启 PortalI18NEditable 时必须配置不同且稳定的编号。
	PortalI18NEditablePbIndex int
	// PortalHiddenLanguages 标识该字段为门户语言隐藏字段。
	// 字段必须是 []string JSON；为空表示所有语言均可展示，命中当前语言时 Portal 查询自动隐藏该记录。
	PortalHiddenLanguages bool

	// 标识 json 数组的查询方式，如传 [1,2] 默认是 or（未设置时默认使用 OR 查询）
	// 另外还支持对 tag 组 [g1:1,g1:2,g2:1,g2:2] 的查询，支持按组 and 或 or 查询，组内的 tag 值 and 或 or 查询
	JSONArrQueryType JSONArrQueryType

	// 索引
	// JSON 数组 CAST 索引 item 的长度 为 0 时不生成索引
	// ALTER TABLE `articles` ADD INDEX category_ids ((CAST(`category_ids` as char(15) ARRAY)))
	IndexJSONArrItemLen int

	// JSON 对象 keys CAST 索引 item 的长度 为 0 时不生成索引
	// ALTER TABLE `articles` ADD INDEX title_i18n ((CAST(JSON_KEYS(`title_i18n`) as char(7) ARRAY)))
	IndexJSONObjKeysLen int

	// PublicReadable 标识该字段可在 /public/ 只读接口中返回（非敏感展示字段）
	// 建议仅用于昵称、头像、名称等可公开展示信息
	PublicReadable bool

	// MyOwner 标识 /my/ 作用域下的归属判断字段（当前用户与数据所有者匹配）
	// 一张表只能配置一个 Owner 字段
	MyOwner bool

	// S3 URL 自动预签名配置
	// 标识该字段的内容中可能包含 S3 私有文件 URL，需要自动处理授权
	// 入库时：清除 URL 中的 S3 授权参数（X-Amz-Signature 等），保留干净的 URL
	// 出库时：自动为匹配的 URL 生成预签名，前端可直接访问
	// 支持纯 URL 字段（如封面 cover）和富文本字段（如 content、markdown）
	// JSON 字段（如 attachments）通过 json.Marshal → 替换 → json.Unmarshal 处理
	S3Presign bool

	// 当字段为结构体时 支持配置 pb 的包名和 import 路径
	PB FieldPB

	// SearchAnalyzer 搜索分词器（如 ik_max_word, ik_smart, standard 等）
	// 仅对 text 类型字段有效；为空则使用默认分词器
	SearchAnalyzer string `json:"searchAnalyzer,omitempty"`

	// SearchType 搜索字段类型覆盖（keyword/text/integer/nested/flat_object 等）
	// 为空则根据 Go 类型自动推断；field.JSON 为结构体/对象切片时默认 flat_object，标量切片按元素类型（如 []string→keyword、[]int64→long）
	SearchType string `json:"searchType,omitempty"`

	// NotSearch 标识该字段不同步到搜索索引
	NotSearch bool `json:"notSearch,omitempty"`
}

// FnHandleType 处理方式
type FnHandleType string

const (
	// TError 返回错误
	TError FnHandleType = "error"
	// TIgnore 忽略
	TIgnore FnHandleType = "ignore"
	// TDefault 默认值
	TDefault FnHandleType = "default"
)

// FieldFunc 字段取值的方法
// 示例
//
//	 &ann.FieldFunc{
//		PkgPath: `"github.com/zhixingdx/public/pkg/middleware/meta"`,
//		Place:   ann.PlaceData,
//		GetStr:  "userID,_:=meta.GetUserID(ctx)",
//		SetStr:  `if b.Creator==nil && userID != "" {b.Creator = &userID}`, service 的写法 b.Field=xxx
//		SetStr: `if b.Creator==nil && userID != "" {db.SetCreator(userID)}`, data 层的写法 db.SetField(xxx)
//	}
type FieldFunc struct {
	PkgPath []string
	// 执行的方法
	GetStr string
	SetStr string
	// 执行的层  TODO biz 层待支持
	Place  Place
	Portal FnPortalType
}

// FnPortalType Portal 层的执行方式
type FnPortalType string

const (
	// FnPortalTypeDefault 默认 都执行
	FnPortalTypeDefault FnPortalType = ""
	// FnPortalTypeOnly 只在 Portal 层执行
	FnPortalTypeOnly FnPortalType = "only"
	// FnPortalTypeNot 不在 Portal 层执行
	FnPortalTypeNot FnPortalType = "not"
)

// PBTimeType pb 的时间类型
type PBTimeType string

const (
	// PBTimeTypeSecond 秒
	PBTimeTypeSecond PBTimeType = "second"
	// PBTimeTypeMillisecond 毫秒
	PBTimeTypeMillisecond PBTimeType = "milli"
)

// JSONArrQueryType JSON 数组查询类型
// 用于标识 JSON 数组字段的查询方式
type JSONArrQueryType string

const (
	// JSONArrQueryTypeDefault 默认值，默认使用 OR 查询（向后兼容）
	JSONArrQueryTypeDefault JSONArrQueryType = ""
	// JSONArrQueryTypeAnd 普通数组 AND 查询，数组中的所有值都必须存在（交集）
	// 示例: [1,2] 表示字段必须同时包含 1 和 2
	JSONArrQueryTypeAnd JSONArrQueryType = "and"
	// JSONArrQueryTypeOr 普通数组 OR 查询，数组中任意一个值存在即可（并集）
	// 示例: [1,2] 表示字段包含 1 或 2 即可
	// 注意: 当 JSONArrQueryType 为空或未设置时，默认使用 OR 查询
	JSONArrQueryTypeOr JSONArrQueryType = "or"
	// JSONArrQueryTypeGroupAndItemOr Tag 组查询，组间 AND + 组内 OR
	// 示例: [g1:1,g1:2,g2:1,g2:2] 表示 (g1组包含1或包含2) AND (g2组包含1或包含2)
	JSONArrQueryTypeGroupAndItemOr JSONArrQueryType = "group_and_item_or"
	// JSONArrQueryTypeGroupOrItemAnd Tag 组查询，组间 OR + 组内 AND
	// 示例: [g1:1,g1:2,g2:1,g2:2] 表示 (g1组包含1且包含2) OR (g2组包含1且包含2)
	JSONArrQueryTypeGroupOrItemAnd JSONArrQueryType = "group_or_item_and"
)

// FieldQuery 字段的查询
type FieldQuery struct {
	// 默认支持 eq 不启用时全部关闭
	Disabled bool
	// 不等于
	NotEq bool
	// 是否开启包含
	In bool
	// 是否开启不包含
	NotIn bool
	// 是否开启模糊查询
	Like bool
	// 是否开启不模糊查询
	NotLike bool
	// 是否开启范围查询
	Range bool
}

// 当字段为结构体时 支持配置 pb 的包名和 import 路径
type FieldPB struct {
	Package string
	Import  string
}

// ObjectStoreMode 对象存储访问模式
type ObjectStoreMode string

const (
	// ObjectStoreModeInline 内联模式：读取时从 S3 下载内容，填充到原字段返回
	// 适用于前端需要编辑的内容（富文本、JSON 配置等）
	ObjectStoreModeInline ObjectStoreMode = "inline"
	// ObjectStoreModeReference 引用模式：读取时生成预签名 URL，返回 URL 和文件大小
	// 适用于只读结果数据（分析报告、AI 生成内容等）
	ObjectStoreModeReference ObjectStoreMode = "reference"
)

// ObjectStoreContentType 对象存储内容类型，决定文件后缀和 S3 Content-Type
type ObjectStoreContentType string

const (
	// ObjectStoreContentTypeAuto 自动推断（按内容首字符判断）
	ObjectStoreContentTypeAuto ObjectStoreContentType = ""
	// ObjectStoreContentTypeJSON application/json
	ObjectStoreContentTypeJSON ObjectStoreContentType = "json"
	// ObjectStoreContentTypeHTML text/html
	ObjectStoreContentTypeHTML ObjectStoreContentType = "html"
	// ObjectStoreContentTypeMarkdown text/markdown
	ObjectStoreContentTypeMarkdown ObjectStoreContentType = "md"
	// ObjectStoreContentTypeText text/plain
	ObjectStoreContentTypeText ObjectStoreContentType = "txt"
)

// ObjectStoreConfig 大字段对象存储配置
type ObjectStoreConfig struct {
	// 存储模式: inline(返回内容) / reference(返回URL)
	Mode ObjectStoreMode `json:"mode"`
	// 超过此大小(字节)时上传到对象存储, 低于则直接存 DB
	// 默认 0 表示始终上传, 设置 > 0 时低于阈值的内容直接存 DB
	Threshold int `json:"threshold"`
	// 仅 inline 模式生效, reference 模式自动忽略, 默认 true
	Compress *bool `json:"compress"`
	// 内容类型, 决定文件后缀和 S3 Content-Type, 空则自动推断
	ContentType ObjectStoreContentType `json:"contentType"`
	// 对象存储目录前缀, 默认 {serviceName}/{tableName}/{fieldName}
	Dir string `json:"dir"`
	// reference 模式: 预签名 URL 有效期(秒), 默认 900 (15分钟)
	PresignExpiry int `json:"presignExpiry"`
	// 是否在列表(Query)接口返回, 默认 false
	ShowInList bool `json:"showInList"`
	// inline 模式: 是否缓存从 S3 下载的内容到 Redis, 默认 false
	InlineCache bool `json:"inlineCache"`
	// inline 模式: 内容缓存过期时间(秒), 默认 600 (10分钟)
	InlineCacheTTL int `json:"inlineCacheTTL"`
	// 删除/编辑记录时是否同步删除 S3 对象, 默认 false
	// S3 存储成本极低, 且直接删除可能误删被其他记录引用的对象
	DeleteOnRemove bool `json:"deleteOnRemove"`
}

// GetPresignExpiry 获取预签名有效期(秒)，默认 900
func (c *ObjectStoreConfig) GetPresignExpiry() int {
	if c.PresignExpiry > 0 {
		return c.PresignExpiry
	}
	return 900
}

// GetInlineCacheTTL 获取内容缓存 TTL(秒)，默认 600
func (c *ObjectStoreConfig) GetInlineCacheTTL() int {
	if c.InlineCacheTTL > 0 {
		return c.InlineCacheTTL
	}
	return 600
}

// IsCompress 是否压缩（仅 inline 模式生效，reference 模式强制不压缩）
func (c *ObjectStoreConfig) IsCompress() bool {
	if c.Mode == ObjectStoreModeReference {
		return false
	}
	if c.Compress != nil {
		return *c.Compress
	}
	return true
}

// IsSearchExcluded 判断字段是否应排除在搜索 API 之外（biz/service/proto 层使用）
// OnlyData 字段不对外暴露，RowIrreversibleEncrypt 不可逆加密无法搜索
// SoftDeleteField 不入搜索索引（软删除记录直接从索引删除，存此字段无意义）
// Encrypt 和 MaskEncrypt 保留：加密字段仍有搜索和展示需求
func (f *Field) IsSearchExcluded() bool {
	return f.OnlyData || f.NotSearch || f.RowIrreversibleEncrypt || f.SoftDeleteField
}

// IsSearchDataExcluded 判断字段是否应排除在搜索索引之外（Data 层使用）
// 仅排除不可逆加密、显式标记不搜索、软删除字段
// OnlyData（如 phoneIndex）需要入库用于加密匹配搜索
func (f *Field) IsSearchDataExcluded() bool {
	return f.NotSearch || f.RowIrreversibleEncrypt || f.SoftDeleteField
}

// FieldNameKey 字段的注解名称
const FieldNameKey = "Field"

// Name of the annotation. Used by the codegen templates.
func (Field) Name() string {
	return FieldNameKey
}

// GetFieldConfig 获取字段的配置
func GetFieldConfig(node *gen.Field) *Field {
	ann := node.Annotations[FieldNameKey]
	if ann == nil {
		return &Field{}
	}
	data, err := json.Marshal(ann)
	if err != nil {
		log.Fatal(err)
	}
	var conf Field
	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal(err)
	}
	return &conf
}

// GetFieldsConfig 获取所有字段的配置
func GetFieldsConfig(node *gen.Type) map[string]*Field {
	fields := make(map[string]*Field)
	indexMap := make(map[int]bool)

	idConfig := GetFieldConfig(node.ID)
	if idConfig.PbIndex == 0 {
		idConfig.PbIndex = 1
	}
	if !isValidProtoFieldNumber(idConfig.PbIndex) {
		log.Fatalf("字段 %s 的 pbIndex %d 不是合法的 Protobuf 字段号", node.ID.Name, idConfig.PbIndex)
	}
	indexMap[idConfig.PbIndex] = true
	fields[node.ID.Name] = idConfig

	for _, field := range node.Fields {
		config := GetFieldConfig(field)
		fields[field.Name] = config
		if config.PbIndex == 0 && !config.OnlyData {
			log.Fatalf("字段 %s 的 pbIndex 未配置", field.Name)
		}
		if config.PbIndex != 0 && !isValidProtoFieldNumber(config.PbIndex) {
			log.Fatalf("字段 %s 的 pbIndex %d 不是合法的 Protobuf 字段号", field.Name, config.PbIndex)
		}

		if config.PbIndex == 1 && !strings.EqualFold(field.Name, "id") {
			log.Fatalf("字段 %s 的 pbIndex 不能等于 1, 默认 1 预留给 ID", field.Name)
		}

		if indexMap[config.PbIndex] && config.PbIndex != 0 && !strings.EqualFold(field.Name, "id") {
			log.Fatalf("字段 %s 的 pbIndex 重复", field.Name)
		}
		indexMap[config.PbIndex] = true
		if config.OnlyUnique {
			_, iB := lo.Find(node.Indexes, func(item *gen.Index) bool {
				return len(item.Columns) == 1 && strings.EqualFold(item.Columns[0], field.Name)
			})
			if !iB {
				log.Fatalf("字段 %s 的 onlyUnique 为 true 时必须有单字段 Unique 索引", field.Name)
			}
		}
		if config.I18NFor != "" {
			// 必须是一个 json 字段
			if !field.IsJSON() {
				log.Fatalf("字段 %s 为多语言字段时时值必须为 json 类型", field.Name)
			}
			// 必须是一个 map 并且 key 为 string
			if !strings.HasPrefix(field.Type.String(), "map[string]") {
				log.Fatalf("字段 %s 为多语言字段时时值必须为 map[string]* 类型", field.Name)
			}
			cField := strings.TrimSuffix(field.Name, "I18n")
			if _, ok := fields[cField]; !ok {
				log.Fatalf("字段 %s 为多语言字段时必须有对应的值字段 %s", field.Name, cField)
			}
		}
		if config.PortalHiddenLanguages {
			if !field.IsJSON() || field.Type.String() != "[]string" {
				log.Fatalf("字段 %s 为门户语言隐藏字段时必须为 []string JSON 类型", field.Name)
			}
		}
	}
	if errMsg := validatePortalI18NEditableFields(fields); errMsg != "" {
		log.Fatal(errMsg)
	}
	return fields
}

// validatePortalI18NEditableFields 校验 Portal 完整多语言字段依赖的本地投影和独立 Proto 编号。
func validatePortalI18NEditableFields(fields map[string]*Field) string {
	fieldNames := make([]string, 0, len(fields))
	for fieldName, config := range fields {
		if config == nil {
			continue
		}
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)

	indexMap := make(map[int]string)
	for _, fieldName := range fieldNames {
		config := fields[fieldName]
		if config.PbIndex > 0 {
			indexMap[config.PbIndex] = fieldName
		}
	}

	portalI18NIndexMap := make(map[int]string)
	for _, fieldName := range fieldNames {
		config := fields[fieldName]
		if !config.PortalI18NEditable {
			continue
		}
		if strings.TrimSpace(config.I18NFor) == "" {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时必须配置 I18NFor", fieldName)
		}
		if !config.PortalOnlyLocalLang {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时必须同时开启 PortalOnlyLocalLang", fieldName)
		}
		if config.NotPortalEdit {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时不能同时开启 NotPortalEdit", fieldName)
		}
		if config.NotPortalDetailReply {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时不能同时开启 NotPortalDetailReply", fieldName)
		}
		// Portal 详情字段集合继承通用详情集合；允许 NotDetailReply 会形成可写却不可读回填的不对称契约。
		if config.NotDetailReply {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时不能同时开启 NotDetailReply", fieldName)
		}
		// OnlyData/OnlyMeta 会同时退出编辑和回复字段集，必须拒绝静默开启但不生成任何读写契约。
		if config.OnlyData {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时不能同时开启 OnlyData", fieldName)
		}
		if config.OnlyMeta {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时不能同时开启 OnlyMeta", fieldName)
		}
		if config.PortalI18NEditablePbIndex <= 0 {
			return fmt.Sprintf("字段 %s 开启 PortalI18NEditable 时必须配置 PortalI18NEditablePbIndex", fieldName)
		}
		if !isValidProtoFieldNumber(config.PortalI18NEditablePbIndex) {
			return fmt.Sprintf("字段 %s 的 PortalI18NEditablePbIndex %d 不是合法的 Protobuf 字段号", fieldName, config.PortalI18NEditablePbIndex)
		}
		// 完整 map 与 local 字段存在于同一个详情消息中，必须提前阻止字段号复用，否则只能等 protoc 才能发现冲突。
		if otherField, exists := indexMap[config.PortalI18NEditablePbIndex]; exists {
			return fmt.Sprintf("字段 %s 的 PortalI18NEditablePbIndex 与字段 %s 的 pbIndex %d 重复", fieldName, otherField, config.PortalI18NEditablePbIndex)
		}
		if otherField, exists := portalI18NIndexMap[config.PortalI18NEditablePbIndex]; exists {
			return fmt.Sprintf("字段 %s 与字段 %s 的 PortalI18NEditablePbIndex %d 重复", fieldName, otherField, config.PortalI18NEditablePbIndex)
		}
		portalI18NIndexMap[config.PortalI18NEditablePbIndex] = fieldName
	}
	return ""
}

// isValidProtoFieldNumber 判断字段号是否位于 Protobuf 可用范围且避开内部保留区间。
func isValidProtoFieldNumber(index int) bool {
	return index > 0 &&
		index <= maxProtoFieldNumber &&
		(index < protoReservedFieldNumberStart || index > protoReservedFieldNumberEnd)
}
