// Package versionedresource 提供版本化资源运行时共享的确定性正文摘要能力。
package versionedresource

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ContentHash 对持久化后的 payload 字段信封做规范化编码并返回小写 SHA-256。
func ContentHash(payload map[string]any) (string, error) {
	normalized := make(map[string]any, len(payload))
	for key, value := range payload {
		if marked, ok := value.(optionalObject); ok {
			raw, err := CanonicalJSON(marked.value)
			if err != nil {
				return "", err
			}
			if string(raw) == "null" || string(raw) == "{}" {
				continue
			}
			value = marked.value
		}
		normalized[key] = value
	}
	canonical, err := CanonicalJSON(normalized)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:]), nil
}

// CanonicalJSON 把支持的 Go 值转换为稳定 JSON；对象键排序，数组顺序和 nil/空值差异均保留。
func CanonicalJSON(value any) ([]byte, error) {
	normalized, err := normalizeCanonicalValue(reflect.ValueOf(value))
	if err != nil {
		return nil, err
	}
	return json.Marshal(normalized)
}

// normalizeCanonicalValue 递归归一化反射值，并拒绝不能跨运行时稳定表达的类型。
func normalizeCanonicalValue(value reflect.Value) (any, error) { // nolint:gocyclo
	if !value.IsValid() {
		return nil, nil
	}
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, nil
		}
		return normalizeCanonicalValue(value.Elem())
	}
	if value.CanInterface() {
		if instant, ok := value.Interface().(time.Time); ok {
			return instant.UTC().Format(time.RFC3339Nano), nil
		}
		if raw, ok := value.Interface().(json.RawMessage); ok {
			return decodeCanonicalJSON(raw)
		}
		if message, ok := value.Interface().(proto.Message); ok {
			raw, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(message)
			if err != nil {
				return nil, fmt.Errorf("编码 protobuf payload 失败: %w", err)
			}
			return decodeCanonicalJSON(raw)
		}
	}
	switch value.Kind() {
	case reflect.Bool:
		return value.Bool(), nil
	case reflect.String:
		return value.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint(), nil
	case reflect.Float32, reflect.Float64:
		floatValue := value.Float()
		if math.IsNaN(floatValue) || math.IsInf(floatValue, 0) {
			return nil, fmt.Errorf("payload 浮点数必须是有限值")
		}
		return floatValue, nil
	case reflect.Slice:
		if value.IsNil() {
			return nil, nil
		}
		if value.Type().Elem().Kind() == reflect.Uint8 {
			// encoding/json 对 []byte 使用标准 base64，显式复制避免后续修改影响结果。
			return append([]byte(nil), value.Bytes()...), nil
		}
		fallthrough
	case reflect.Array:
		items := make([]any, value.Len())
		for index := 0; index < value.Len(); index++ {
			item, err := normalizeCanonicalValue(value.Index(index))
			if err != nil {
				return nil, err
			}
			items[index] = item
		}
		return items, nil
	case reflect.Map:
		if value.IsNil() {
			return nil, nil
		}
		if value.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("payload map 只支持 string 键")
		}
		result := make(map[string]any, value.Len())
		iterator := value.MapRange()
		for iterator.Next() {
			item, err := normalizeCanonicalValue(iterator.Value())
			if err != nil {
				return nil, err
			}
			result[iterator.Key().String()] = item
		}
		return result, nil
	case reflect.Struct:
		if !value.CanInterface() {
			return nil, fmt.Errorf("payload struct %s 不可访问", value.Type())
		}
		raw, err := json.Marshal(value.Interface())
		if err != nil {
			return nil, fmt.Errorf("编码 payload struct %s 失败: %w", value.Type(), err)
		}
		return decodeCanonicalJSON(raw)
	default:
		return nil, fmt.Errorf("payload 不支持类型 %s", value.Type())
	}
}

// decodeCanonicalJSON 使用 json.Number 保留整数精度，再进入统一递归归一化。
func decodeCanonicalJSON(raw []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("解析 payload JSON 失败: %w", err)
	}
	return value, nil
}

// optionalObject 标记显式允许忽略缺省 JSON 对象的单个版本字段。
type optionalObject struct{ value any }

// OmitEmptyObject 只在正文摘要中省略 nil 或空对象；数组和标量保持原值。
func OmitEmptyObject(value any) any { return optionalObject{value: value} }
