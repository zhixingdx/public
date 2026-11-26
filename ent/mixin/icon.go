package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Icon mixin 通常用于表示图标图片地址
type Icon struct {
	mixin.Schema
}

// Fields _
func (Icon) Fields() []ent.Field {
	return []ent.Field{
		field.String("icon").
			MaxLen(1023).
			Optional().
			Default("").
			Comment("图标").
			Annotations(ann.Field{
				PbIndex:        236,
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
func (Icon) Index() []ent.Index {
	return []ent.Index{}
}
