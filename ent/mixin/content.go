// Package mixin switch
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Content 定义通用大文本内容，默认按 HTML 富文本执行 XSS 过滤。
type Content struct {
	mixin.Schema

	// DisableXSSFilter 仅供 JSON 等不得被 HTML 清洗改写的非 HTML 内容，启用后由业务层负责协议校验与输出安全。
	DisableXSSFilter bool
}

// Fields 返回内容字段，空配置始终保持 XSS 过滤开启的安全默认。
func (m Content) Fields() []ent.Field {
	return []ent.Field{
		field.Text("content").
			Optional().
			Default("").
			Comment("内容").
			Annotations(ann.Field{
				PbIndex:   216,
				S3Presign: true,
				XSSFilter: !m.DisableXSSFilter,
				Query: ann.FieldQuery{
					Disabled: true,
				},
				NotQueryReply:       true,
				NotPortalQueryReply: true,
			}),
	}
}

// Index _
func (Content) Index() []ent.Index {
	return []ent.Index{}
}
