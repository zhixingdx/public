package enterr

import (
	"regexp"
	"strings"
)

const (
	// mySQLErrDuplicateEntry MySQL 驱动常见唯一冲突片段（如 Error 1062）。
	mySQLErrDuplicateEntry = "Duplicate entry"
	// pgsqlErrDuplicateKey PostgreSQL 唯一冲突错误串常见片段。
	pgsqlErrDuplicateKey = "duplicate key value"
)

// MySQL Duplicate entry 'value' for key 'table.idx_name' 中连续单引号段（非空捕获）。
var mysqlDuplicateQuotedSegments = regexp.MustCompile(`'([^']+)'`)

// PostgreSQL Detail: Key (column)=(value) already exists.
var pgDuplicateKeyDetail = regexp.MustCompile(`Key \(([^)]+)\)=\(([^)]+)\)`)

// IsMySQLUniqueViolationMessage 判断错误文案是否像 MySQL 唯一约束冲突。
func IsMySQLUniqueViolationMessage(errMsg string) bool {
	return strings.Contains(errMsg, mySQLErrDuplicateEntry)
}

// IsPostgreSQLUniqueViolationMessage 判断错误文案是否像 PostgreSQL 唯一约束冲突。
func IsPostgreSQLUniqueViolationMessage(errMsg string) bool {
	return strings.Contains(errMsg, pgsqlErrDuplicateKey)
}

// ParseMySQLUniqueViolation 解析 MySQL 唯一索引冲突：metadata 键取自索引名按下划线分段后的第二段（与历史 EntHandler 行为一致）。
// 解析失败返回 (nil, false)，调用方应仍返回 duplicateKey 多语言提示但不带 metadata。
func ParseMySQLUniqueViolation(errMsg string) (map[string]string, bool) {
	matches := mysqlDuplicateQuotedSegments.FindAllStringSubmatch(errMsg, -1)
	if len(matches) < 2 {
		return nil, false
	}
	value := matches[0][1]
	if value == "" {
		return nil, false
	}
	fieldParts := strings.Split(matches[1][1], ".")
	if len(fieldParts) != 2 {
		return nil, false
	}
	indexName := fieldParts[1]
	indexSegs := strings.Split(indexName, "_")
	if len(indexSegs) < 2 {
		return nil, false
	}
	return map[string]string{indexSegs[1]: value}, true
}

// ParsePostgreSQLUniqueViolation 从错误详情中解析 Key (col)=(val)，成功则返回 map[column]value。
func ParsePostgreSQLUniqueViolation(errMsg string) (map[string]string, bool) {
	matches := pgDuplicateKeyDetail.FindStringSubmatch(errMsg)
	if len(matches) < 3 {
		return nil, false
	}
	column, value := matches[1], matches[2]
	if column == "" || value == "" {
		return nil, false
	}
	return map[string]string{column: value}, true
}
