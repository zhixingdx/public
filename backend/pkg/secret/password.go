package secret

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

const passwordRounds = 600000

// HashPassword 每次设置口令采用独立随机盐，保存算法、成本及校验值。
func HashPassword(password string) (string, error) {
	if len(password) == 0 || len(password) > 1024 {
		return "", errors.New("口令长度无效")
	}
	var salt [16]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return "", err
	}
	key, err := pbkdf2.Key(sha256.New, password, salt[:], passwordRounds, 32)
	if err != nil {
		return "", err
	}
	enc := base64.RawStdEncoding
	return strings.Join([]string{"pbkdf2-sha256", strconv.Itoa(passwordRounds), enc.EncodeToString(salt[:]), enc.EncodeToString(key)}, "$"), nil
}

// VerifyPassword 拒绝未知格式和越界成本，比较使用常量时间。
func VerifyPassword(encoded, password string) bool {
	if encoded == "" || len(encoded) > 255 || len(password) == 0 || len(password) > 1024 {
		return false
	}
	p := strings.Split(encoded, "$")
	if len(p) != 4 || p[0] != "pbkdf2-sha256" {
		return false
	}
	rounds, err := strconv.Atoi(p[1])
	if err != nil || rounds < passwordRounds || rounds > 2000000 {
		return false
	}
	enc := base64.RawStdEncoding.Strict()
	salt, err := enc.DecodeString(p[2])
	if err != nil || len(salt) != 16 {
		return false
	}
	want, err := enc.DecodeString(p[3])
	if err != nil || len(want) != 32 {
		return false
	}
	got, err := pbkdf2.Key(sha256.New, password, salt, rounds, len(want))
	return err == nil && subtle.ConstantTimeCompare(got, want) == 1
}
