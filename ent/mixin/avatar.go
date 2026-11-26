package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Avatar mixin 通常用于表示 头像 图片地址
type Avatar struct {
	mixin.Schema
}

// Fields _
func (Avatar) Fields() []ent.Field {
	return []ent.Field{
		field.String("avatar").
			MaxLen(1023).
			Optional().
			Default("").
			Comment("头像").
			Annotations(ann.Field{
				PbIndex:        233,
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
func (Avatar) Index() []ent.Index {
	return []ent.Index{}
}
