// Package mixin mail
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
	"github.com/zhixingdx/oipivot/backend/pkg/utils"
)

// Mail mixin 邮箱
type Mail struct {
	mixin.Schema
}

// Fields _
func (Mail) Fields() []ent.Field {
	return []ent.Field{
		field.String("mail").
			Comment("邮箱").
			MaxLen(63).
			Optional().
			Annotations(
				ann.Field{
					PbIndex:     209,
					Validate:    "{ max_len:63, email:true }",
					MaskEncrypt: utils.MaskTypeEmail,
					Query: ann.FieldQuery{
						NotEq: true,
						In:    true,
						NotIn: true,
					},
				},
			),

		field.String("mailIndex").MaxLen(64).Optional().Sensitive().Annotations(ann.Field{OnlyData: true}),

		field.String("mailCipher").
			MaxLen(255).
			Optional().
			Sensitive().
			Annotations(
				ann.Field{
					OnlyData:  true,
					NotSearch: true,
				},
			),
	}
}

// Index _
func (Mail) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("mailIndex"),
	}
}
