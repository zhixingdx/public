// Package bizcall 校验跨服务业务引用。
package bizcall

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/zhixingdx/public/api/apicommon/v1"
)

func Valid(call *apicommon.BizCall) bool {
	return validField(call.GetService()) && validField(call.GetAction()) && (call.GetBizId() == "" || validField(call.GetBizId()))
}

func validField(value string) bool {
	return len(value) > 0 && len(value) <= 127 && utf8.ValidString(value) && strings.TrimSpace(value) == value && strings.IndexFunc(value, unicode.IsControl) < 0
}
