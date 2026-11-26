// Package meta student
package meta

import (
	"context"

	"github.com/go-kratos/kratos/v3/errors"
)

// GetStudentID 获取学生 id。
func GetStudentID(ctx context.Context) (string, *errors.Error) {
	return GetValue(ctx, studentIDKey)
}

// SetStudentID 设置学生 id。
func SetStudentID(ctx context.Context, studentID string) context.Context {
	return SetValue(ctx, studentIDKey, studentID)
}
