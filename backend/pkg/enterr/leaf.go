package enterr

import (
	"errors"
	"strings"
)

// LeafMessage 沿标准库 errors.Unwrap 链遍历，取最内层错误的 Error() 文本（去首尾空白）。
// 用于 Ent validator 包装链末端文案（如 "value is greater than the required length"）。
func LeafMessage(err error) string {
	if err == nil {
		return ""
	}
	var leaf string
	for e := err; e != nil; e = errors.Unwrap(e) {
		leaf = strings.TrimSpace(e.Error())
	}
	return leaf
}
