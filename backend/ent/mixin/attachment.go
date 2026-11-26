// Package mixin switch
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot-api/apicommon/v1"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
)

const portalAttachmentI18NPbIndex = 315

// Attachment  mixin 通常用于表示 简介	存放 html 之类的文本信息
type Attachment struct {
	mixin.Schema
}

// Fields _
func (Attachment) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("attachments", []*apicommon.FileAttachment{}).Optional().
			Default(make([]*apicommon.FileAttachment, 0)).
			Comment("附件").
			Annotations(ann.Field{
				PbIndex:   223,
				S3Presign: true,
				Query:     ann.FieldQuery{Disabled: true},
				PB: ann.FieldPB{
					Package: "api.apicommon",
					Import:  "apicommon/v1/file.proto",
				},
			}),
	}
}

// Index 返回附件字段使用的索引定义。
func (Attachment) Index() []ent.Index {
	return []ent.Index{}
}

// AttachmentI18n mixin 通常用于表示 国际化附件
type AttachmentI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造国际化附件字段，并把实例级 Portal 编辑能力写入字段注解。
func (m AttachmentI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("attachmentsI18n", map[string][]*apicommon.FileAttachment{}).Optional().Default(make(map[string][]*apicommon.FileAttachment)).Comment("国际化附件").
			Annotations(ann.Field{
				PbIndex:                   307,
				S3Presign:                 true,
				I18NFor:                   "attachments",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalAttachmentI18NPbIndex,
				Query:                     ann.FieldQuery{Disabled: true},
				PB: ann.FieldPB{
					Package: "api.apicommon",
					Import:  "apicommon/v1/file.proto",
				},
			}),
	}
}

// Attachment  mixin 通常用于表示 简介	存放 html 之类的文本信息
type AttachmentNotQueryReply struct {
	mixin.Schema
}

// Fields _
func (AttachmentNotQueryReply) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("attachments", []*apicommon.FileAttachment{}).Optional().
			Default(make([]*apicommon.FileAttachment, 0)).
			Comment("附件").
			Annotations(ann.Field{
				PbIndex:       223,
				S3Presign:     true,
				NotQueryReply: true,
				Query:         ann.FieldQuery{Disabled: true},
				PB: ann.FieldPB{
					Package: "api.apicommon",
					Import:  "apicommon/v1/file.proto",
				},
			}),
	}
}

// Index 返回不进入查询响应的附件字段所使用的索引定义。
func (AttachmentNotQueryReply) Index() []ent.Index {
	return []ent.Index{}
}

// AttachmentI18n mixin 通常用于表示 国际化附件
type AttachmentI18nNotQueryReply struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造不进入查询列表的国际化附件字段，并把实例级 Portal 编辑能力写入字段注解。
func (m AttachmentI18nNotQueryReply) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("attachmentsI18n", map[string][]*apicommon.FileAttachment{}).Optional().Default(make(map[string][]*apicommon.FileAttachment)).Comment("国际化附件").
			Annotations(ann.Field{
				PbIndex:                   307,
				S3Presign:                 true,
				I18NFor:                   "attachments",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalAttachmentI18NPbIndex,
				NotQueryReply:             true,
				Query:                     ann.FieldQuery{Disabled: true},
				PB: ann.FieldPB{
					Package: "api.apicommon",
					Import:  "apicommon/v1/file.proto",
				},
			}),
	}
}
