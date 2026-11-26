// Package mixin biz_call
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
)

// BizCall mixin 内部调用的相关字段
type BizCall struct {
	mixin.Schema
	// RequireSource 要求每条记录都有来源服务和业务动作。
	RequireSource bool
}

// Fields _
func (b BizCall) Fields() []ent.Field {
	service := field.String("service").MaxLen(127)
	action := field.String("action").MaxLen(127)
	if b.RequireSource {
		service.NotEmpty().Immutable()
		action.NotEmpty().Immutable()
	} else {
		service.Optional().Default("")
		action.Optional().Default("")
	}
	return []ent.Field{
		service.
			Comment("服务").
			Annotations(ann.Field{
				PbIndex:              227,
				Validate:             "{max_len:127}",
				NotPortalAdd:         true,
				NotPortalEdit:        true,
				NotPortalQuery:       true,
				NotPortalQueryReply:  true,
				NotPortalDetailReply: true,
			}),
		action.
			Comment("操作").
			Annotations(ann.Field{
				PbIndex:              228,
				Validate:             "{max_len:127}",
				NotPortalAdd:         true,
				NotPortalEdit:        true,
				NotPortalQuery:       true,
				NotPortalQueryReply:  true,
				NotPortalDetailReply: true,
			}),
		// 业务关联 ID,
		field.String("bizID").
			MaxLen(127).
			Optional().
			Default("").
			Comment("业务 ID").
			Annotations(ann.Field{
				PbIndex:       229,
				Validate:      "{max_len:127}",
				NotPortalAdd:  true,
				NotPortalEdit: true,
				Query: ann.FieldQuery{
					In:    true,
					NotIn: true,
				},
			}),
	}
}
