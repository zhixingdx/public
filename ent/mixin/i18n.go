// Package mixin i18n 多语言字段
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	portalExtraI18NPbIndex    = 308
	portalContentI18NPbIndex  = 309
	portalSummaryI18NPbIndex  = 310
	portalCoverI18NPbIndex    = 311
	portalTitleI18NPbIndex    = 312
	portalSubTitleI18NPbIndex = 313
	portalMarkdownI18NPbIndex = 314
)

// ExtraI18n 提供国际化扩展信息字段，并可按实例开放 Portal 完整编辑能力。
type ExtraI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// HiddenLanguages 资源级语言隐藏配置。
// 为空表示所有语言均可展示；Portal 查询会按当前语言上下文自动排除命中的内容。
type HiddenLanguages struct {
	mixin.Schema
}

// Fields _
func (HiddenLanguages) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("hiddenLanguages", []string{}).
			Optional().
			Default([]string{}).
			Comment("隐藏语言").
			Annotations(ann.Field{
				PbIndex:               299,
				PortalHiddenLanguages: true,
				IndexJSONArrItemLen:   15,
				NotPortalAdd:          true,
				NotPortalEdit:         true,
				NotPortalQuery:        true,
				NotPortalQueryReply:   true,
				NotPortalDetailReply:  true,
				Query:                 ann.FieldQuery{Disabled: true},
			}),
	}
}

// Fields 构造国际化扩展信息字段，并把实例级 Portal 编辑能力写入字段注解。
func (m ExtraI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("extraI18n", map[string]*structpb.Struct{}).
			Optional().
			Default(map[string]*structpb.Struct{}).
			Comment("国际化扩展信息").
			Annotations(ann.Field{
				PbIndex:                   300,
				I18NFor:                   "extra",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalExtraI18NPbIndex,
				Query:                     ann.FieldQuery{Disabled: true},
			}),
	}
}

// Content  mixin 通常用于表示 内容	存放 html 之类的大文本信息
type ContentI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造国际化内容字段，并把实例级 Portal 编辑能力写入字段注解。
func (m ContentI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("contentI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化内容").
			Annotations(ann.Field{
				PbIndex:                   301,
				I18NFor:                   "content",
				S3Presign:                 true,
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalContentI18NPbIndex,
				XSSFilter:                 true,
				Query:                     ann.FieldQuery{Disabled: true},
				// 查询不返回 只在详情中返回
				NotQueryReply:       true,
				NotPortalQueryReply: true,
			}),
	}
}

// SummaryI18n mixin 通常用于表示 简介
type SummaryI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造国际化简介字段，并把实例级 Portal 编辑能力写入字段注解。
func (m SummaryI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("summaryI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化简介").
			Annotations(ann.Field{
				PbIndex:                   302,
				I18NFor:                   "summary",
				S3Presign:                 true,
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalSummaryI18NPbIndex,
				XSSFilter:                 true,
				Query:                     ann.FieldQuery{Disabled: true},
			}),
	}
}

// Cover  mixin 通常用于表示 封面 图片地址
type CoverI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造国际化封面字段，并把实例级 Portal 编辑能力写入字段注解。
func (m CoverI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("coverI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化封面").
			Annotations(ann.Field{
				PbIndex:                   303,
				I18NFor:                   "cover",
				S3Presign:                 true,
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalCoverI18NPbIndex,
				Query:                     ann.FieldQuery{Disabled: true},
			}),
	}
}

// TitleI18n  mixin 通常用于表示 标题
type TitleI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造国际化标题字段，并把实例级 Portal 编辑能力写入字段注解。
func (m TitleI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("titleI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化标题").
			Annotations(ann.Field{
				PbIndex:                   304,
				I18NFor:                   "title",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalTitleI18NPbIndex,
				IndexJSONObjKeysLen:       7,
				Query: ann.FieldQuery{
					Like:    true,
					NotLike: true,
				},
			}),
	}
}

// 副标题国际化
type SubTitleI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造国际化副标题字段，并把实例级 Portal 编辑能力写入字段注解。
func (m SubTitleI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("subTitleI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化副标题").
			Annotations(ann.Field{
				PbIndex:                   305,
				I18NFor:                   "subTitle",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalSubTitleI18NPbIndex,
				Query:                     ann.FieldQuery{Disabled: true},
			}),
	}
}

// TitleI18n  mixin 通常用于表示 标题
type TitlePublicReadableI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造公开可读的国际化标题字段，并把实例级 Portal 编辑能力写入字段注解。
func (m TitlePublicReadableI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("titleI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化标题").
			Annotations(ann.Field{
				PbIndex:                   304,
				I18NFor:                   "title",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalTitleI18NPbIndex,
				IndexJSONObjKeysLen:       7,
				Query: ann.FieldQuery{
					Like:    true,
					NotLike: true,
				},
				PublicReadable: true,
			}),
	}
}

// 副标题国际化
type SubTitlePublicReadableI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造公开可读的国际化副标题字段，并把实例级 Portal 编辑能力写入字段注解。
func (m SubTitlePublicReadableI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("subTitleI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化副标题").
			Annotations(ann.Field{
				PbIndex:                   305,
				I18NFor:                   "subTitle",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalSubTitleI18NPbIndex,
				Query:                     ann.FieldQuery{Disabled: true},
				PublicReadable:            true,
			}),
	}
}

// MarkdownI18n mixin 通常用于表示 内容 markdown
type MarkdownI18n struct {
	mixin.Schema
	// PortalI18NEditable 显式允许 Portal 在详情和编辑接口中管理完整多语言内容。
	PortalI18NEditable bool
}

// Fields 构造国际化 Markdown 字段，并把实例级 Portal 编辑能力写入字段注解。
func (m MarkdownI18n) Fields() []ent.Field {
	return []ent.Field{
		field.JSON("markdownI18n", map[string]string{}).
			Optional().
			Default(make(map[string]string)).
			Comment("国际化内容 markdown").
			Annotations(ann.Field{
				PbIndex:                   306,
				I18NFor:                   "markdown",
				PortalOnlyLocalLang:       true,
				PortalI18NEditable:        m.PortalI18NEditable,
				PortalI18NEditablePbIndex: portalMarkdownI18NPbIndex,
				Query:                     ann.FieldQuery{Disabled: true},
				NotQueryReply:             true,
				NotPortalQueryReply:       true,
			}),
	}
}
