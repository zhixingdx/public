package enterr

import (
	"errors"
	"reflect"
	"strings"
)

// joinErrorsUnwrap 与标准库 errors.Join 生成的包装错误一致：Unwrap 返回多段子错误。
type joinErrorsUnwrap interface {
	Unwrap() []error
}

// collectValidationFieldLeaves 遍历错误树，收集 Ent 生成的 ValidationError 节点（见 isEntValidationErrorLike）。
// 键为 schema 字段名（须通过 IsSafeSchemaFieldName），值为该节点上的 LeafMessage（unwrap 最内层文案）。
// errors.Join 多段校验时同一字段后出现覆盖先出现的项。
func collectValidationFieldLeaves(err error) map[string]string {
	out := make(map[string]string)
	if err == nil {
		return out
	}
	var walk func(error)
	walk = func(e error) {
		if e == nil {
			return
		}
		if je, ok := e.(joinErrorsUnwrap); ok {
			for _, sub := range je.Unwrap() {
				walk(sub)
			}
			return
		}
		if isEntValidationErrorLike(e) {
			name := strings.TrimSpace(entValidationErrorName(e))
			if IsSafeSchemaFieldName(name) {
				out[name] = LeafMessage(e)
			}
		}
		if u := errors.Unwrap(e); u != nil {
			walk(u)
		}
	}
	walk(err)
	return out
}

// isEntValidationErrorLike 用反射识别 ent 生成代码中的 ValidationError：结构体名固定为 ValidationError，且含可导出 string 字段 Name。
// 依赖 ent 代码生成约定；若未来 ent 改名需同步调整。
func isEntValidationErrorLike(e error) bool {
	if e == nil {
		return false
	}
	t := reflect.TypeOf(e)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t == nil || t.Kind() != reflect.Struct || t.Name() != "ValidationError" {
		return false
	}
	nf, ok := t.FieldByName("Name")
	return ok && nf.Type.Kind() == reflect.String
}

func entValidationErrorName(e error) string {
	v := reflect.ValueOf(e)
	if v.Kind() == reflect.Pointer && !v.IsNil() {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ""
	}
	f := v.FieldByName("Name")
	if !f.IsValid() || f.Kind() != reflect.String || !f.CanInterface() {
		return ""
	}
	s, _ := f.Interface().(string)
	return s
}
