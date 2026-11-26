package enterr

import (
	"regexp"
	"strings"
)

// schemaFieldNameRe 与 ent Schema 字段常见命名一致（小驼峰），用于限制 metadata 键，防止异常键名。
var schemaFieldNameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

// IsSafeSchemaFieldName 判断名称是否可作为对前端/API 暴露的 metadata 键。
func IsSafeSchemaFieldName(name string) bool {
	name = strings.TrimSpace(name)
	return name != "" && schemaFieldNameRe.MatchString(name)
}
