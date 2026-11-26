// Package mixin switch
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Title  mixin 通常用于表示 标题
type TitlePublicReadable struct {
	mixin.Schema
}

// Fields _
func (TitlePublicReadable) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			MaxLen(511).
			Comment("标题").
			Annotations(ann.Field{
				PbIndex:  219,
				Validate: "{max_len:511}",
				Query: ann.FieldQuery{
					Like:    true,
					NotLike: true,
					In:      true,
				},
				PublicReadable: true,
			}),
	}
}

// Index _
func (TitlePublicReadable) Index() []ent.Index {
	return []ent.Index{
		index.Fields("title"),
	}
}

// 副标题
type SubTitlePublicReadable struct {
	mixin.Schema
}

// Fields _
func (SubTitlePublicReadable) Fields() []ent.Field {
	return []ent.Field{
		field.String("subTitle").
			MaxLen(511).
			Optional().
			Default("").
			Comment("副标题").
			Annotations(ann.Field{
				PbIndex:  222,
				Validate: "{max_len:511}",
				Query: ann.FieldQuery{
					Disabled: true,
				},
				PublicReadable: true,
			}),
	}
}
