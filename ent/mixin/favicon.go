package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Favicon mixin 通常用于表示 Favicon 图标地址
type Favicon struct {
	mixin.Schema
}

// Fields _
func (Favicon) Fields() []ent.Field {
	return []ent.Field{
		field.String("favicon").
			MaxLen(1023).
			Optional().
			Default("").
			Comment("favicon").
			Annotations(ann.Field{
				PbIndex:        235,
				S3Presign:      true,
				PublicReadable: true,
				Validate:       "{max_len:1023}",
				Query: ann.FieldQuery{
					Disabled: true,
				},
			}),
	}
}

// Index _
func (Favicon) Index() []ent.Index {
	return []ent.Index{}
}
