// Package secret 为生成的数据层加密和读取部署凭据。
package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"strings"
)

const prefix = "gcm1:"

func aead(key string) (cipher.AEAD, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	return cipher.NewGCMWithRandomNonce(block)
}

// Encrypt 返回带格式标记和随机 nonce 的认证密文。
func Encrypt(value, key string) (string, error) {
	a, err := aead(key)
	if err != nil {
		return "", err
	}
	return prefix + base64.RawStdEncoding.EncodeToString(a.Seal(nil, nil, []byte(value), nil)), nil
}

// Decrypt 检查完整性后返回部署凭据。
func Decrypt(value, key string) (string, error) {
	body, ok := strings.CutPrefix(value, prefix)
	if !ok {
		return "", errors.New("未知凭据格式")
	}
	raw, err := base64.RawStdEncoding.DecodeString(body)
	if err != nil {
		return "", err
	}
	a, err := aead(key)
	if err != nil {
		return "", err
	}
	p, err := a.Open(nil, nil, raw, nil)
	return string(p), err
}
