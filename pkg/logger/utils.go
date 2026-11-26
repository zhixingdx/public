// Package logger utils.go
package logger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// toString 将日志值安全转换为字符串，typed nil 指针不调用其 String 方法。
func toString(v interface{}) string {
	var key string
	if v == nil {
		return key
	}
	// error/Stringer接口可能承载nil指针；日志不应掩盖原业务错误或使成功请求panic。
	if value := reflect.ValueOf(v); value.Kind() == reflect.Ptr && value.IsNil() {
		return ""
	}
	switch v := v.(type) {
	case float64:
		key = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		key = strconv.FormatFloat(float64(v), 'f', -1, 32)
	case int:
		key = strconv.Itoa(v)
	case uint:
		key = strconv.FormatUint(uint64(v), 10)
	case int8:
		key = strconv.Itoa(int(v))
	case uint8:
		key = strconv.FormatUint(uint64(v), 10)
	case int16:
		key = strconv.Itoa(int(v))
	case uint16:
		key = strconv.FormatUint(uint64(v), 10)
	case int32:
		key = strconv.Itoa(int(v))
	case uint32:
		key = strconv.FormatUint(uint64(v), 10)
	case int64:
		key = strconv.FormatInt(v, 10)
	case uint64:
		key = strconv.FormatUint(v, 10)
	case string:
		key = v
	case bool:
		key = strconv.FormatBool(v)
	case []byte:
		key = string(v)
	case fmt.Stringer:
		key = v.String()
	default:
		newValue, _ := json.Marshal(v)
		key = string(newValue)
	}
	return key
}
