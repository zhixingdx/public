// Package mixin password
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Password mixin 密码
type Password struct {
	mixin.Schema
}

// Fields _
func (Password) Fields() []ent.Field {
	return []ent.Field{
		field.String("password").
			Sensitive().
			MaxLen(255).
			Optional().
			Default("").
			Comment("密码").
			Annotations(
				ann.Field{
					PbIndex:                210,
					Validate:               "{max_len:255}",
					RowIrreversibleEncrypt: true,
					NotQueryReply:          true,
					NotDetailReply:         true,
					NotPortalQuery:         true,
					NotPortalQueryReply:    true,
					NotPortalDetailReply:   true,
					Query: ann.FieldQuery{
						Disabled: true,
					},
				},
			),
	}
}
