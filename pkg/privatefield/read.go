package privatefield

import (
	"context"
	"sync"

	"entgo.io/ent/dialect/sql"
)

// Reject 使条件构造失败时无法退化成无条件读取或修改。
func Reject(err error) func(*sql.Selector) {
	return func(s *sql.Selector) { s.Where(sql.False()); s.AddError(err) }
}

type readKey struct{}
type readState struct {
	sync.Mutex
	err error
}

// Read 在调用方已授权的用例内启用明文投影。任何字段打开失败都会放弃整个结果。
func Read[T any](ctx context.Context, operation func(context.Context) (T, error)) (T, error) {
	state := new(readState)
	value, err := operation(context.WithValue(ctx, readKey{}, state))
	state.Lock()
	defer state.Unlock()
	if state.err != nil {
		var zero T
		return zero, state.err
	}
	return value, err
}
func Reading(ctx context.Context) bool { return ctx.Value(readKey{}) != nil }
func ReadError(ctx context.Context, err error) {
	if state, ok := ctx.Value(readKey{}).(*readState); ok && err != nil {
		state.Lock()
		defer state.Unlock()
		state.err = ErrProtection
	}
}
