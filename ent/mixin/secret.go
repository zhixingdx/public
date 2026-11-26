// Package mixin secret
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Secret ent
type Secret struct {
	mixin.Schema
	// OptionalID 允许同表中的其他渠道使用独立认证方式。
	OptionalID bool
}

// Fields Secret
func (s Secret) Fields() []ent.Field {
	id := field.String("secretID")
	if s.OptionalID {
		id.Optional()
	}
	return []ent.Field{
		id.
			Comment("密钥ID").
			MaxLen(127).
			Annotations(ann.Field{
				PbIndex:  211,
				Validate: "{max_len:127}",
				Query: ann.FieldQuery{
					NotEq: true,
					In:    true,
					NotIn: true,
				},
			}),

		field.String("secretKey").
			Optional().
			MaxLen(512).
			Default("").
			Comment("密钥Key").
			Annotations(ann.Field{
				PbIndex:        212,
				Validate:       "{max_len:255}",
				Encrypt:        true,
				NotQueryReply:  true,
				NotDetailReply: true,
			},
			),
	}
}

// Indexes PhoneMxin
func (Secret) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("secretID"),
	}
}
