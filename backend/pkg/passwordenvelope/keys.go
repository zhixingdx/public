// Package passwordenvelope 打开门户端 JWE。挑战的签发和单次消费归 base-safe。
package passwordenvelope

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"
	"strings"

	jose "github.com/go-jose/go-jose/v4"
)

var ErrInvalid = errors.New("口令信封无效")

type Keys struct {
	active string
	keys   map[string]*ecdsa.PrivateKey
}
type Payload struct {
	Challenge string `json:"challenge"`
	Value     string `json:"value"`
}

// Load 接受 PKCS#8 PEM 密钥束。第一把签发公钥，其余仅用于打开轮换前的信封。
func Load(path string) (*Keys, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := &Keys{keys: map[string]*ecdsa.PrivateKey{}}
	for len(strings.TrimSpace(string(data))) > 0 {
		block, rest := pem.Decode(data)
		if block == nil || block.Type != "PRIVATE KEY" {
			return nil, errors.New("需要 PKCS#8 P-256 私钥")
		}
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.New("口令信封密钥格式无效")
		}
		ec, ok := key.(*ecdsa.PrivateKey)
		if !ok || ec.Curve != elliptic.P256() {
			return nil, errors.New("口令信封要求 P-256 私钥")
		}
		jwk := jose.JSONWebKey{Key: &ec.PublicKey}
		thumb, err := jwk.Thumbprint(crypto.SHA256)
		if err != nil {
			return nil, err
		}
		kid := base64.RawURLEncoding.EncodeToString(thumb)
		if out.active == "" {
			out.active = kid
		}
		out.keys[kid] = ec
		data = rest
	}
	if out.active == "" {
		return nil, errors.New("缺少口令信封密钥")
	}
	return out, nil
}
func (k *Keys) Public() (kid, x, y string) {
	key := k.keys[k.active]
	enc := base64.RawURLEncoding
	return k.active, enc.EncodeToString(key.X.FillBytes(make([]byte, 32))), enc.EncodeToString(key.Y.FillBytes(make([]byte, 32)))
}

func (k *Keys) Open(compact string) (Payload, error) {
	if len(compact) > 12288 {
		return Payload{}, ErrInvalid
	}
	obj, err := jose.ParseEncryptedCompact(compact, []jose.KeyAlgorithm{jose.ECDH_ES}, []jose.ContentEncryption{jose.A256GCM})
	if err != nil {
		return Payload{}, ErrInvalid
	}
	if _, ok := obj.Header.ExtraHeaders["zip"]; ok {
		return Payload{}, ErrInvalid
	}
	key := k.keys[obj.Header.KeyID]
	if key == nil {
		return Payload{}, ErrInvalid
	}
	plain, err := obj.Decrypt(key)
	if err != nil {
		return Payload{}, ErrInvalid
	}
	var p Payload
	if json.Unmarshal(plain, &p) != nil || len(p.Challenge) != 43 || len(p.Value) == 0 || len(p.Value) > 1024 {
		return Payload{}, ErrInvalid
	}
	return p, nil
}
