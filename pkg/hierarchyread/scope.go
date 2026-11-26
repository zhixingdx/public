// Package hierarchyread 保存进程内已授权的下级读取谓词，不改变请求身份或序列化为出站 metadata。
package hierarchyread

import (
	"context"

	"entgo.io/ent/dialect/sql"
)

// PredicateFactory 绑定静态业务资源，并由生成的实体租户列构造关联 SQL 条件。
type PredicateFactory func(tenantColumn string) func(*sql.Selector)

// scopeKey 使用包私有类型隔离能力，外部 metadata 不能伪造读取授权。
type scopeKey struct{}

// readScope 只携带授权后建立的资源及谓词工厂，不记录可展开的下级列表。
type readScope struct {
	resource string
	factory  PredicateFactory
}

// WithScope 由授权 owner 附加资源专属读取能力，保留原租户、Realm 和员工上下文。
func WithScope(ctx context.Context, resource string, factory PredicateFactory) context.Context {
	return context.WithValue(ctx, scopeKey{}, readScope{resource: resource, factory: factory})
}

// Predicate 仅在当前实体静态资源匹配时返回范围，不允许同一个授权扩大其他表的读取。
func Predicate(ctx context.Context, resource, tenantColumn string) (func(*sql.Selector), bool) {
	scope, ok := ctx.Value(scopeKey{}).(readScope)
	if !ok || resource == "" || scope.resource != resource || tenantColumn == "" {
		return nil, false
	}
	if scope.factory == nil {
		return denyScope, true
	}
	predicate := scope.factory(tenantColumn)
	if predicate == nil {
		return denyScope, true
	}
	return predicate, true
}

// denyScope 让异常能力失败关闭，不能因工厂缺失退回不受下级范围约束的查询。
func denyScope(selector *sql.Selector) {
	selector.Where(sql.False())
}
