// Package external 解码第三方供应商规定的密码报文。
package external

import (
	"crypto/aes"
	"encoding/base64"
	"errors"
)

const QdeduProfile = "qdedu-profile"

var errInvalid = errors.New("第三方加密报文无效")

// Decode 由供应商适配器固定协议。青岛 OAuth 用户资料采用 Base64、AES-ECB 和 PKCS#7。
func Decode(protocol, payload string, key []byte) (string, error) {
	if protocol != QdeduProfile {
		return "", errInvalid
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errInvalid
	}
	raw, err := base64.StdEncoding.DecodeString(payload)
	if err != nil || len(raw) == 0 || len(raw)%aes.BlockSize != 0 {
		return "", errInvalid
	}
	for i := 0; i < len(raw); i += aes.BlockSize {
		block.Decrypt(raw[i:i+aes.BlockSize], raw[i:i+aes.BlockSize])
	}
	padding := int(raw[len(raw)-1])
	if padding < 1 || padding > aes.BlockSize {
		return "", errInvalid
	}
	for _, b := range raw[len(raw)-padding:] {
		if int(b) != padding {
			return "", errInvalid
		}
	}
	return string(raw[:len(raw)-padding]), nil
}
