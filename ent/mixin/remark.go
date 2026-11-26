// Package mixin mail
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Remark mixin
type Remark struct {
	mixin.Schema
}

// Fields _
func (Remark) Fields() []ent.Field {
	return []ent.Field{
		field.String("remark").
			Comment("备注").
			MaxLen(1023).
			Optional().
			Default("").
			Annotations(
				ann.Field{
					PbIndex:  214,
					Validate: "{max_len:1023}",
					Query: ann.FieldQuery{
						Like: true,
					},
					SearchType:           "text",
					NotPortalAdd:         true,
					NotPortalEdit:        true,
					NotPortalQuery:       true,
					NotPortalQueryReply:  true,
					NotPortalDetailReply: true,
				},
			),
	}
}
