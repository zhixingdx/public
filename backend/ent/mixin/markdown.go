// Package mixin switch
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
)

// Markdown  mixin 通常用于表示 内容 markdown
type Markdown struct {
	mixin.Schema
}

// Fields _
func (Markdown) Fields() []ent.Field {
	return []ent.Field{
		field.Text("markdown").Optional().Default("").Comment("内容 markdown").Annotations(
			ann.Field{
				PbIndex:             230,
				S3Presign:           true,
				Query:               ann.FieldQuery{Disabled: true},
				NotQueryReply:       true,
				NotPortalQueryReply: true,
			}),
	}
}

// Index _
func (Markdown) Index() []ent.Index {
	return []ent.Index{}
}
