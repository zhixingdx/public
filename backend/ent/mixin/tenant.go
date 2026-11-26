// Package mixin appid
package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/zhixingdx/oipivot/backend/ent/ann"
	"github.com/zhixingdx/oipivot/backend/pkg/dict"
)

// TenantID mixin 通常用于 sass 场景中标识不同的 租户
type TenantID struct {
	mixin.Schema
}

// Fields 定义由生成器参与租户隔离的不可变归属字段。
func (TenantID) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenantID").
			MaxLen(63).
			Comment("租户 ID").
			Immutable().
			Annotations(ann.Field{
				PbIndex:   205,
				SassField: true,
				Validate:  "{max_len:63}",
			}),
	}
}

// Index 为租户归属查询声明索引。
func (TenantID) Index() []ent.Index {
	return []ent.Index{
		index.Fields("tenantID"),
	}
}

// TenantIDOnlyPortal mixin 通常用于 sass 场景中标识不同的 租户
type TenantIDOnlyPortal struct {
	mixin.Schema
}

// Fields 定义仅在 Portal 侧自动参与租户隔离的不可变归属字段。
func (TenantIDOnlyPortal) Fields() []ent.Field {
	return []ent.Field{
		field.String("tenantID").
			MaxLen(63).
			Comment("租户 ID").
			Immutable().
			Annotations(ann.Field{
				PbIndex:        205,
				SassField:      true,
				SassOnlyPortal: true,
				Validate:       "{max_len:63}",
			}),
	}
}

// Index 为门户租户归属查询声明索引。
func (TenantIDOnlyPortal) Index() []ent.Index {
	return []ent.Index{
		index.Fields("tenantID"),
	}
}

// TenantVisibilityScope 定义租户可见性；个人范围须显式启用且配合唯一 MyOwner 字段。
type TenantVisibilityScope struct {
	mixin.Schema
	AllowSelf bool
}

// Fields 定义租户可见范围、可见租户与不可变 Realm 字段。
func (m TenantVisibilityScope) Fields() []ent.Field {
	scope := field.Enum("visibilityScope").Comment("租户可见范围").Optional().GoType(dict.TenantVisibilityScope(dict.TenantVisibilityScopeInternal)).Default(dict.TenantVisibilityScopeInternal.String()).Annotations(ann.Field{PbIndex: 231, Query: ann.FieldQuery{In: true, NotIn: true}})
	if m.AllowSelf {
		// 保留未接入表的三种枚举和既有生成结果，防止缺少本人鉴权的表提前接受 self。
		scope.Values(dict.TenantVisibilityScopeSelf.String())
	}
	return []ent.Field{
		scope,
		field.JSON("visibleTenants", []string{}).Optional().Default([]string{}).Comment("可见租户").Annotations(ann.Field{PbIndex: 232, IndexJSONArrItemLen: 15}),
		field.String("realm").
			MinLen(2).
			MaxLen(31).
			Default("default").
			Immutable().
			Comment("所属业务渠道").
			Annotations(ann.Field{
				PbIndex:    244,
				OnlyData:   true,
				AutoCreate: true,
				DefaultFn: &ann.FieldFunc{
					PkgPath: []string{"github.com/zhixingdx/oipivot/backend/pkg/middleware/meta"},
					Place:   ann.PlaceData,
					GetStr: `currentRealm, realmErr := meta.GetRealm(ctx)
if realmErr != nil { return nil, realmErr }`,
					SetStr: `db.SetRealm(currentRealm)`,
				},
			}),
	}
}

// Indexes 为 Realm 与可见范围组合查询提供 Ent 可识别的索引。
func (TenantVisibilityScope) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("realm", "visibilityScope"),
	}
}
