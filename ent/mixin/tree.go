// Package mixin tree helpers.
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
)

const (
	treeRootIDPbIndex = 237
	treeSortPbIndex   = 238
)

// TreeRootID 为邻接表树提供根节点 ID 快速过滤字段。
type TreeRootID struct {
	mixin.Schema
}

// Fields 返回树根字段定义。
func (TreeRootID) Fields() []ent.Field {
	return []ent.Field{
		field.String("rootID").
			MaxLen(31).
			Optional().
			Comment("树根ID").
			Annotations(ann.Field{
				PbIndex:    treeRootIDPbIndex,
				TreeRootID: true,
				Validate:   "{max_len:31}",
				Query: ann.FieldQuery{
					In:    true,
					NotIn: true,
				},
			}),
	}
}

// Indexes 返回树根索引。
func (TreeRootID) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("rootID"),
	}
}

// TreeSort 为树节点提供同级排序字段。
type TreeSort struct {
	mixin.Schema
}

// Fields of the TreeSort.
func (TreeSort) Fields() []ent.Field {
	return []ent.Field{
		field.Int("sort").Optional().Default(0).Comment("排序").Annotations(
			ann.Field{
				PbIndex:  treeSortPbIndex,
				TreeSort: true,
			},
		),
	}
}
