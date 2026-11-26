package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Logo mixin 通常用于表示 Logo 图片地址
type Logo struct {
	mixin.Schema
}

// Fields _
func (Logo) Fields() []ent.Field {
	return []ent.Field{
		field.String("logo").
			MaxLen(1023).
			Optional().
			Default("").
			Comment("logo").
			Annotations(ann.Field{
				PbIndex:        234,
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
func (Logo) Index() []ent.Index {
	return []ent.Index{}
}
