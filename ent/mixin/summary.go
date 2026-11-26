// Package mixin switch
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/public/ent/ann"
)

// Summary  mixin 通常用于表示 简介	存放 html 之类的文本信息
type Summary struct {
	mixin.Schema
}

// Fields _
func (Summary) Fields() []ent.Field {
	return []ent.Field{
		field.Text("summary").Optional().
			Default("").
			Comment("简介").
			Annotations(ann.Field{
				PbIndex:    217,
				S3Presign:  true,
				XSSFilter:  true,
				SearchType: "text",
				Query: ann.FieldQuery{
					Disabled: true,
				},
			}),
	}
}

// Index _
func (Summary) Index() []ent.Index {
	return []ent.Index{}
}
