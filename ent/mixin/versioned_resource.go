// Package mixin 提供版本化资源的公共 Ent 字段。
package mixin

import (
	"regexp"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
	"github.com/zhixingdx/public/pkg/dict"
)

// versionedResourceContentHashPattern 固定内部 SHA-256 为小写十六进制规范编码。
var versionedResourceContentHashPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

const (
	// VersionedResourceRevisionPBIndex 是版本修订号的官方 mixin 协议编号。
	VersionedResourceRevisionPBIndex = 239
	// VersionedResourceStatusPBIndex 是版本生命周期的官方 mixin 协议编号。
	VersionedResourceStatusPBIndex = 240
	// VersionedResourceContentHashPBIndex 是版本内容摘要的官方 mixin 协议编号。
	VersionedResourceContentHashPBIndex = 241
	// VersionedResourcePublishedTimePBIndex 是版本发布时间的官方 mixin 协议编号。
	VersionedResourcePublishedTimePBIndex = 242
	// VersionedResourcePublishedByPBIndex 由互斥的员工或用户发布者字段共用。
	VersionedResourcePublishedByPBIndex = 243
)

// VersionedResourceVersion 提供所有版本内容表共享的生命周期固定字段。
type VersionedResourceVersion struct {
	mixin.Schema
}

// Fields 返回公共版本字段，业务表不能通过普通 CRUD 改写。
func (VersionedResourceVersion) Fields() []ent.Field {
	return versionedResourceVersionFields()
}

// VersionedResourcePublishedByStaff 提供员工发布者字段。
type VersionedResourcePublishedByStaff struct {
	mixin.Schema
}

// Fields 返回员工发布者字段；同一版本表只能选择一种发布者 mixin。
func (VersionedResourcePublishedByStaff) Fields() []ent.Field {
	return []ent.Field{publishedByField("publishedByStaff", "发布员工 ID")}
}

// VersionedResourcePublishedByUser 提供终端用户发布者字段。
type VersionedResourcePublishedByUser struct {
	mixin.Schema
}

// Fields 返回用户发布者字段，与员工发布者字段共用协议编号。
func (VersionedResourcePublishedByUser) Fields() []ent.Field {
	return []ent.Field{publishedByField("publishedByUser", "发布用户 ID")}
}

// versionedResourceVersionFields 构造两个发布者变体共享的版本控制字段。
func versionedResourceVersionFields() []ent.Field {
	return []ent.Field{
		field.Int64("revision").
			Positive().
			Immutable().
			Comment("资源内单调递增修订号").
			Annotations(ann.Field{
				PbIndex:    VersionedResourceRevisionPBIndex,
				AutoCreate: true,
				AutoUpdate: true,
				Query:      ann.FieldQuery{In: true, NotIn: true, Range: true},
			}),
		field.Enum("status").
			GoType(dict.VersionStatus(dict.VersionStatusDraft)).
			Default(string(dict.VersionStatusDraft)).
			Comment("版本生命周期").
			Annotations(ann.Field{
				PbIndex:    VersionedResourceStatusPBIndex,
				AutoCreate: true,
				AutoUpdate: true,
				Query:      ann.FieldQuery{In: true, NotIn: true},
			}),
		field.String("contentHash").
			MaxLen(64).
			Match(versionedResourceContentHashPattern).
			SchemaType(map[string]string{
				dialect.MySQL: "char(64) CHARACTER SET ascii COLLATE ascii_bin",
			}).
			Comment("完整 payload 规范化 SHA-256").
			Annotations(ann.Field{
				PbIndex:    VersionedResourceContentHashPBIndex,
				Validate:   "{len:64}",
				AutoCreate: true,
				AutoUpdate: true,
				Query:      ann.FieldQuery{In: true},
			}),
		field.Time("publishedTime").
			Optional().
			Nillable().
			Comment("发布时间").
			SchemaType(map[string]string{"mysql": "datetime"}).
			Annotations(ann.Field{
				PbIndex:      VersionedResourcePublishedTimePBIndex,
				AutoCreate:   true,
				AutoUpdate:   true,
				PBTimeToType: ann.PBTimeTypeSecond,
				Query:        ann.FieldQuery{Range: true},
			}),
	}
}

// publishedByField 统一发布者字段约束，避免员工与用户变体产生协议和校验漂移。
func publishedByField(name, comment string) ent.Field {
	return field.String(name).
		MaxLen(15).
		Optional().
		Default("").
		Comment(comment).
		Annotations(ann.Field{
			// 发布者身份只用于管理与审计，匿名门户不得返回员工或用户标识。
			NotPortalQueryReply:  true,
			NotPortalDetailReply: true,
			PbIndex:              VersionedResourcePublishedByPBIndex,
			Validate:             "{max_len:15}",
			AutoCreate:           true,
			AutoUpdate:           true,
		})
}
