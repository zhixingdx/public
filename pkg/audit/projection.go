package audit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	pb "github.com/zhixingdx/oipivot-api/base/audit/v1"
	"github.com/zhixingdx/oipivot/backend/pkg/versionedresource"
)

// ErrProjection 表示源值不能安全、完整地编码，仅供内部诊断及边界本地化。
var ErrProjection = errors.New("invalid audit projection")

// Field 固定生成时批准的物理列和编码方式，运行时不接受请求决定字段范围。
type Field struct {
	Name   string
	Column string
	Kind   string
	System bool
}

// Values 用字段存在性保存 SQL NULL；JSON null 仍以原始 null 文本存在。
type Values map[string]json.RawMessage

// ValueText 读取事件身份列的无损文本，字符串去引号而整数保持十进制。
func ValueText(value json.RawMessage) string {
	var text string
	if json.Unmarshal(value, &text) == nil {
		return text
	}
	return string(value)
}

// CredentialName 为已分类凭据提供额外名称兜底，不能用于推断普通字段已加密。
func CredentialName(name string) bool {
	name = strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(name))
	switch name {
	case "password", "passwd", "passwordnonce", "passwordcipher", "secret", "secretkey", "clientsecret", "secretaccesskey", "accesskey", "authorization", "cookie", "token", "accesstoken", "refreshtoken", "privatekey", "verificationcode", "captcha", "passwordhash", "passwordsalt", "apikey", "otp":
		return true
	}
	return false
}

// Pointer 把 Schema 字段名编码为单个 RFC6901 路径段。
func Pointer(name string) string {
	return "/" + strings.ReplaceAll(strings.ReplaceAll(name, "~", "~0"), "/", "~1")
}

// CanonicalValue 在源端拒绝重复键、凭据子树及深度越界，并无损规范化数字。
func CanonicalValue(raw []byte) (json.RawMessage, error) {
	if len(raw) == 0 || !utf8.Valid(raw) || len(raw) > 64*1024 {
		return nil, ErrProjection
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := checkJSONValue(decoder, 0); err != nil {
		return nil, err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return nil, ErrProjection
	}
	value, err := versionedresource.CanonicalJSON(json.RawMessage(raw))
	if err != nil || len(value) > 64*1024 {
		return nil, ErrProjection
	}
	return value, nil
}

// checkJSONValue 在解码覆盖重复键之前检查结构和动态凭据。
func checkJSONValue(decoder *json.Decoder, depth int) error {
	token, err := decoder.Token()
	if err != nil {
		return ErrProjection
	}
	delim, compound := token.(json.Delim)
	if !compound {
		return nil
	}
	if depth >= 16 {
		return ErrProjection
	}
	seen := map[string]bool{}
	for decoder.More() {
		if delim == '{' {
			key, err := decoder.Token()
			if err != nil {
				return ErrProjection
			}
			name, ok := key.(string)
			if !ok || seen[name] || CredentialName(name) {
				return ErrProjection
			}
			seen[name] = true
		}
		if err := checkJSONValue(decoder, depth+1); err != nil {
			return err
		}
	}
	end, err := decoder.Token()
	if err != nil || (delim == '{' && end != json.Delim('}')) || (delim == '[' && end != json.Delim(']')) {
		return ErrProjection
	}
	return nil
}

// SQLValue 直接编码驱动返回的原值；nil 单独交由调用者表示 SQL NULL。
func SQLValue(kind string, value any) (json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}
	if raw, ok := value.([]byte); ok {
		value = string(raw)
	}
	if kind == "json" {
		raw, ok := value.(string)
		if !ok {
			return nil, ErrProjection
		}
		return CanonicalValue([]byte(raw))
	}
	if kind == "number" {
		return CanonicalValue([]byte(fmt.Sprint(value)))
	}
	if kind == "bool" {
		switch fmt.Sprint(value) {
		case "1", "true":
			return json.RawMessage("true"), nil
		case "0", "false":
			return json.RawMessage("false"), nil
		default:
			return nil, ErrProjection
		}
	}
	if instant, ok := value.(time.Time); ok {
		value = instant.UTC()
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, ErrProjection
	}
	return CanonicalValue(raw)
}

// EntityValue 编码 Ent 新建实体字段；JSON 直接使用实体值，避免批量创建后逐行回读数据库。
func EntityValue(kind string, value any) (json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}
	reflected := reflect.ValueOf(value)
	for reflected.Kind() == reflect.Interface || reflected.Kind() == reflect.Pointer {
		if reflected.IsNil() {
			return nil, nil
		}
		reflected = reflected.Elem()
	}
	value = reflected.Interface()
	if kind != "json" {
		return SQLValue(kind, value)
	}
	var raw []byte
	switch typed := value.(type) {
	case json.RawMessage:
		raw = typed
	case []byte:
		raw = typed
	case string:
		raw = []byte(typed)
	default:
		encoded, err := json.Marshal(value)
		if err != nil {
			return nil, ErrProjection
		}
		raw = encoded
	}
	return CanonicalValue(raw)
}

// Snapshot 冻结完整应采字段，不用省略字段掩盖编码失败。
func Snapshot(values Values) (string, error) {
	raw, err := json.Marshal(values)
	if err != nil {
		return "", ErrProjection
	}
	canonical, err := CanonicalValue(raw)
	return string(canonical), err
}

// Changes 按真实存储前后值计算差异，自动更新时间不能独自产生业务事件。
func Changes(fields []Field, before, after Values) ([]*pb.Change, error) {
	changes := []*pb.Change{}
	business := false
	for _, field := range fields {
		old, oldPresent := before[field.Name]
		next, nextPresent := after[field.Name]
		if oldPresent == nextPresent && bytes.Equal(old, next) {
			continue
		}
		change := &pb.Change{Path: Pointer(field.Name), BeforePresent: oldPresent, AfterPresent: nextPresent}
		if oldPresent {
			text := string(old)
			change.Before = &text
		}
		if nextPresent {
			text := string(next)
			change.After = &text
		}
		changes = append(changes, change)
		business = business || !field.System
	}
	if !business {
		return nil, nil
	}
	if len(changes) > 200 {
		return nil, ErrProjection
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	raw, err := json.Marshal(changes)
	if err != nil || len(raw) > 64*1024 {
		return nil, ErrProjection
	}
	return changes, nil
}
