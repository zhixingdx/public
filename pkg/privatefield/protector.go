// Package privatefield 保护可恢复且需要精确查询的个人信息。
package privatefield

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Scope 使用资源的实际归属。同值在不同租户、表或字段中不可关联。
type Scope struct {
	Tenant string
	Field  string
}

// Protector 由数据层装配，生产密码服务可实现同一接口且无需导出密钥。
type Protector interface {
	Seal(context.Context, Scope, string) (string, error)
	Open(context.Context, Scope, string) (string, error)
	Index(context.Context, Scope, string) (string, error)
}

// Software 供本地联调验证字段行为。加密版本可并存，索引密钥单独轮换。
type Software struct {
	active   string
	keys     map[string]cipher.AEAD
	indexKey []byte
}

var ErrProtection = errors.New("个人信息字段保护失败")

// NewSoftware 接收部署注入的独立随机密钥。业务方不可把口令或 data.secret 当成密钥。
func NewSoftware(active string, keys map[string][]byte, indexKey []byte) (*Software, error) {
	if len(indexKey) != 32 || len(keys[active]) != 32 {
		return nil, ErrProtection
	}
	p := &Software{active: active, keys: make(map[string]cipher.AEAD, len(keys)), indexKey: append([]byte(nil), indexKey...)}
	seen := map[string]bool{string(indexKey): true}
	for version, key := range keys {
		if version == "" || strings.Contains(version, ".") || len(version) > 32 || len(key) != 32 || seen[string(key)] {
			return nil, ErrProtection
		}
		seen[string(key)] = true
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}
		p.keys[version], err = cipher.NewGCMWithRandomNonce(block)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

// 编码数组消除分隔符碰撞，同一密码用途的输入域固定。
func encode(scope Scope, value string) ([]byte, error) {
	if scope.Field == "" || scope.Tenant == "" {
		return nil, ErrProtection
	}
	return json.Marshal([3]string{scope.Tenant, scope.Field, value})
}

func (p *Software) Seal(ctx context.Context, scope Scope, plain string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	aad, err := encode(scope, p.active)
	if err != nil {
		return "", err
	}
	sealed := p.keys[p.active].Seal(nil, nil, []byte(plain), aad)
	return p.active + "." + base64.RawURLEncoding.EncodeToString(sealed), nil
}

func (p *Software) Open(ctx context.Context, scope Scope, value string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	version, body, ok := strings.Cut(value, ".")
	a := p.keys[version]
	if !ok || a == nil {
		return "", ErrProtection
	}
	raw, err := base64.RawURLEncoding.Strict().DecodeString(body)
	if err != nil || len(raw) < a.Overhead() {
		return "", ErrProtection
	}
	aad, err := encode(scope, version)
	if err != nil {
		return "", err
	}
	plain, err := a.Open(nil, nil, raw, aad)
	if err != nil {
		return "", ErrProtection
	}
	return string(plain), nil
}

func (p *Software) Index(ctx context.Context, scope Scope, plain string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	msg, err := encode(scope, plain)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, p.indexKey)
	mac.Write(msg)
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// Protect 在写数据库前生成两列，任一步失败均放弃本次写入。
func Protect(ctx context.Context, p Protector, scope Scope, value string) (string, string, error) {
	if p == nil {
		return "", "", ErrProtection
	}
	index, err := p.Index(ctx, scope, value)
	if err != nil {
		return "", "", err
	}
	sealed, err := p.Seal(ctx, scope, value)
	if err != nil {
		return "", "", err
	}
	return sealed, index, nil
}

// Load 读取部署挂载的字段密钥文件，配置与数据库备份分开交付。
func Load(path string) (*Software, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取字段密钥文件: %w", err)
	}
	var keys struct {
		Active     string            `json:"active"`
		Encryption map[string][]byte `json:"encryption"`
		Index      []byte            `json:"index"`
	}
	if err = json.Unmarshal(raw, &keys); err != nil {
		return nil, ErrProtection
	}
	return NewSoftware(keys.Active, keys.Encryption, keys.Index)
}
