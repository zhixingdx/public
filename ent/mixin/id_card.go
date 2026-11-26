// Package mixin 提供身份证字段复用能力。
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
	"github.com/zhixingdx/oipivot/backend/pkg/utils"
)

// IDCard 为实体提供可空的身份证字段，并由生成层统一维护脱敏值和密文。
type IDCard struct {
	mixin.Schema
}

// Fields 定义对外语义简单的 idCard 字段及仅供数据层使用的密文字段。
func (IDCard) Fields() []ent.Field {
	return []ent.Field{
		field.String("idCard").
			Comment("身份证号").
			MaxLen(18).
			Optional().
			Annotations(
				ann.Field{
					PbIndex:     249,
					MaskEncrypt: utils.MaskTypeIDCard,
					Validate:    "{max_len:18}",
					Query: ann.FieldQuery{
						NotEq: true,
						In:    true,
						NotIn: true,
					},
				},
			),

		field.String("idCardIndex").MaxLen(64).Optional().Sensitive().Annotations(ann.Field{OnlyData: true}),

		field.String("idCardCipher").
			MaxLen(255).
			Sensitive().
			Optional().
			Annotations(
				ann.Field{
					OnlyData:  true,
					NotSearch: true,
				},
			),
	}
}

func (IDCard) Indexes() []ent.Index {
	return []ent.Index{index.Fields("idCardIndex")}
}
