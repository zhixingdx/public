// Package identity 校验项目内服务调用的人员身份。
package identity

import (
	"context"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v3/metadata"
)

// ValidID 校验项目及发起者的 UTF-8 标识。
func ValidID(s string) bool {
	return len(s) > 0 && len(s) <= 63 && utf8.ValidString(s) && strings.TrimSpace(s) == s && strings.IndexFunc(s, unicode.IsControl) < 0
}

type Actor struct {
	Type string
	ID   string
}

// ErrPersonRequired 表示调用方未提供可归属到个人的身份。
var ErrPersonRequired = errors.New("valid staff or user identity required")

// FromContext 仅接受可信调用方转发的人员身份。同一员工兼具用户身份时使用员工身份。
func FromContext(ctx context.Context) (Actor, error) {
	md, _ := metadata.FromServerContext(ctx)
	var actor Actor
	for _, kind := range []string{"user", "staff"} {
		ids := md.Values("x-md-global-" + kind + "-id")
		if len(ids) == 0 {
			continue
		}
		if len(ids) != 1 || !ValidID(ids[0]) {
			return Actor{}, ErrPersonRequired
		}
		actor = Actor{Type: kind, ID: ids[0]}
	}
	if actor.ID == "" {
		return Actor{}, ErrPersonRequired
	}
	return actor, nil
}
