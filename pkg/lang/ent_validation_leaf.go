package lang

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// ent/schema/field 内置校验的叶子错误文案及可兜底透传的自定义英文短语，与 pkg/enterr.LeafMessage 输出的最内层串对齐。
var entValidationDetailSafeASCIIRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9\s.,:'/_-]{0,239}$`)

// entValidationDetailKind Ent ORM 校验叶子多语言路由（不导出，仅供 GetEntValidationFieldDetailMsg）。
type entValidationDetailKind int

const (
	entValidationDetailUnsupported entValidationDetailKind = iota
	entValidationDetailMissingField
	entValidationDetailMissingEdge
	entValidationDetailInvalidQueryField
	entValidationDetailMaxLen
	entValidationDetailMinLen
	entValidationDetailMaxRuneLen
	entValidationDetailMinRuneLen
	entValidationDetailOutOfRange
	entValidationDetailFormatMismatch
	entValidationDetailPassthrough
)

// classifyEntValidationLeaf 将 ORM 校验最内层错误文案分类，供 Common 绑定 MsgKey。
func classifyEntValidationLeaf(leaf string) (kind entValidationDetailKind, safeDetail string) {
	if leaf == "" {
		return entValidationDetailUnsupported, ""
	}
	low := strings.ToLower(leaf)

	switch {
	case strings.Contains(low, "missing required field"):
		return entValidationDetailMissingField, ""
	case strings.Contains(low, "missing required edge"):
		return entValidationDetailMissingEdge, ""
	case strings.Contains(low, "invalid field"):
		return entValidationDetailInvalidQueryField, ""
	case leaf == "value is greater than the required length":
		return entValidationDetailMaxLen, ""
	case leaf == "value is less than the required length":
		return entValidationDetailMinLen, ""
	case leaf == "value is greater than the required rune length":
		return entValidationDetailMaxRuneLen, ""
	case leaf == "value is less than the required rune length":
		return entValidationDetailMinRuneLen, ""
	case leaf == "value out of range":
		return entValidationDetailOutOfRange, ""
	case leaf == "value does not match validation":
		return entValidationDetailFormatMismatch, ""
	}

	if !entValidationDetailSafeASCIIRe.MatchString(leaf) {
		return entValidationDetailUnsupported, ""
	}
	if entValidationDetailRisky(low) {
		return entValidationDetailUnsupported, ""
	}
	return entValidationDetailPassthrough, truncateEntValidationRunes(leaf, 240)
}

func entValidationDetailRisky(low string) bool {
	risky := []string{
		"select ", "insert ", "update ", "delete ", "drop ", "union ",
		"../", "..\\", "/etc/", "password=", "secret=", " bearer ",
		"stack trace", "goroutine ", "0x",
	}
	for _, s := range risky {
		if strings.Contains(low, s) {
			return true
		}
	}
	return false
}

func truncateEntValidationRunes(s string, max int) string {
	if max <= 0 || utf8.RuneCountInString(s) <= max {
		return s
	}
	r := []rune(s)
	return string(r[:max])
}
