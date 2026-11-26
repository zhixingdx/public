// Package mixin 提供实体专属元数据表的公共 Ent 存储合同。
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	// ResourceMetadataSpacePBIndex 是 metadata space 的官方 mixin 协议编号。
	ResourceMetadataSpacePBIndex = 245
	// ResourceMetadataNamePBIndex 是 metadata name 的官方 mixin 协议编号。
	ResourceMetadataNamePBIndex = 246
	// ResourceMetadataValuePBIndex 是 metadata value 的官方 mixin 协议编号。
	ResourceMetadataValuePBIndex = 247
	// ResourceMetadataDocumentPBIndex 是 metadata document 的官方 mixin 协议编号。
	ResourceMetadataDocumentPBIndex = 248

	// ResourceMetadataUniqueIndexName 是 owner 内同键同值唯一索引的稳定名称。
	ResourceMetadataUniqueIndexName = "uq_resource_metadata_owner_space_name_value"
	// ResourceMetadataReverseIndexName 是从可检索值反查 owner 的稳定索引名称。
	ResourceMetadataReverseIndexName = "idx_resource_metadata_reverse"
)

// ResourceMetadata 提供实体专属 metadata 表共享的字段、索引和存储归属注解。
// 字符集、排序规则和方言能力由数据库与服务自己的迁移配置统一决定。
type ResourceMetadata struct {
	mixin.Schema
	// OwnerField 指向具体 metadata Schema 自行声明的 owner ID 字段。
	OwnerField string
	// ScopeField 指向 owner 已存在的可选 SassField；无隔离 owner 留空。
	ScopeField string
}

// Fields 返回行业无关的 space/name/value/document 固定字段。
func (ResourceMetadata) Fields() []ent.Field {
	disabledQuery := ann.FieldQuery{Disabled: true}
	return []ent.Field{
		field.String("space").
			NotEmpty().
			MaxLen(63).
			Comment("行业或业务元数据空间").
			Annotations(ann.Field{
				PbIndex:   ResourceMetadataSpacePBIndex,
				Validate:  "{max_len:63}",
				Query:     disabledQuery,
				NotSearch: true,
			}),
		field.String("name").
			NotEmpty().
			MaxLen(63).
			Comment("空间内元数据属性名").
			Annotations(ann.Field{
				PbIndex:   ResourceMetadataNamePBIndex,
				Validate:  "{max_len:63}",
				Query:     disabledQuery,
				NotSearch: true,
			}),
		field.String("value").
			MaxLen(255).
			Default("").
			Comment("可检索的规范 ID 或 code，空串仅作为 document 行哨兵").
			Annotations(ann.Field{
				PbIndex:   ResourceMetadataValuePBIndex,
				Validate:  "{max_len:255}",
				Query:     disabledQuery,
				NotSearch: true,
			}),
		field.JSON("document", &structpb.Struct{}).
			Optional().
			Comment("不可检索的 JSON object 扩展数据").
			Annotations(ann.Field{
				PbIndex:   ResourceMetadataDocumentPBIndex,
				Query:     disabledQuery,
				NotSearch: true,
			}),
	}
}

// Indexes 返回 owner 内唯一索引和可检索值反向索引。
func (config ResourceMetadata) Indexes() []ent.Index {
	uniqueFields := make([]string, 0, 5)
	reverseFields := make([]string, 0, 5)
	if config.ScopeField != "" {
		uniqueFields = append(uniqueFields, config.ScopeField)
		reverseFields = append(reverseFields, config.ScopeField)
	}
	uniqueFields = append(uniqueFields, config.OwnerField, "space", "name", "value")
	reverseFields = append(reverseFields, "space", "name", "value", config.OwnerField)
	return []ent.Index{
		index.Fields(uniqueFields...).Unique().StorageKey(ResourceMetadataUniqueIndexName),
		index.Fields(reverseFields...).StorageKey(ResourceMetadataReverseIndexName),
	}
}

// Annotations 保存 owner/scope 单一事实源，供 Yimoka 在生成前校验表关系。
func (config ResourceMetadata) Annotations() []schema.Annotation {
	return []schema.Annotation{
		ann.ResourceMetadataStorage{
			OwnerField: config.OwnerField,
			ScopeField: config.ScopeField,
		},
	}
}
