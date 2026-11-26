package enterr

import (
	"errors"

	"github.com/zhixingdx/oipivot-api/fault"
)

// ErrDuplicateKey 供同一次用例在并发首次创建时重新读取，不改变对外错误码。
var ErrDuplicateKey = errors.New("duplicate database key")

// constraintAsDuplicateKeyFault 判断 errMsg 是否为 MySQL/PostgreSQL 唯一冲突语境；是则返回 duplicateKey 语义的 400（解析成功则附加 metadata）。
func constraintAsDuplicateKeyFault(errMsg, duplicateKeyMsg string) (e error, handledAsDuplicate bool) {
	if !IsMySQLUniqueViolationMessage(errMsg) && !IsPostgreSQLUniqueViolationMessage(errMsg) {
		return nil, false
	}
	dup := fault.ErrorBadRequest("%s", duplicateKeyMsg).WithCause(ErrDuplicateKey)
	if IsMySQLUniqueViolationMessage(errMsg) {
		if meta, ok := ParseMySQLUniqueViolation(errMsg); ok {
			return dup.WithMetadata(meta), true
		}
		return dup, true
	}
	if meta, ok := ParsePostgreSQLUniqueViolation(errMsg); ok {
		return dup.WithMetadata(meta), true
	}
	return dup, true
}

// ConstraintFault 将数据库约束错误串转为 fault：属唯一冲突类时返回 duplicateKeyMsg（尽力附带 metadata），否则返回 constraintMsg。
func ConstraintFault(errMsg string, duplicateKeyMsg, constraintMsg string) error {
	if e, ok := constraintAsDuplicateKeyFault(errMsg, duplicateKeyMsg); ok {
		return e
	}
	return fault.ErrorBadRequest("%s", constraintMsg)
}

// ValidationFault 单字段场景的便捷方法（无包装错误树时使用）。
func ValidationFault(fieldName string, validationSummaryMsg, detailMsg string) error {
	baseErr := fault.ErrorBadRequest("%s", validationSummaryMsg)
	if !IsSafeSchemaFieldName(fieldName) || detailMsg == "" {
		return baseErr
	}
	return baseErr.WithMetadata(map[string]string{fieldName: detailMsg})
}

// ValidationFaultFromErr 从错误树收集若干 ValidationError 字段，将每字段的叶子串交给 detailForLeaf 转成对客户端安全文案后写入 metadata。
// detailForLeaf 通常闭包捕获 ctx 与 lang.Common；为 nil 时等价于仅返回主文案。
// 若未解析到合法字段或全部 detail 为空，则不附加 metadata。
func ValidationFaultFromErr(err error, validationSummaryMsg string, detailForLeaf func(leaf string) string) error {
	baseErr := fault.ErrorBadRequest("%s", validationSummaryMsg)
	if err == nil || detailForLeaf == nil {
		return baseErr
	}
	fields := collectValidationFieldLeaves(err)
	if len(fields) == 0 {
		return baseErr
	}
	meta := make(map[string]string, len(fields))
	for field, leaf := range fields {
		if d := detailForLeaf(leaf); d != "" {
			meta[field] = d
		}
	}
	if len(meta) == 0 {
		return baseErr
	}
	return baseErr.WithMetadata(meta)
}
