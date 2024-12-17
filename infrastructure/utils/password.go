package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateSalt 生成随机盐
func GenerateSalt() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return string(Base58Encode(bytes)), nil
}

// Encode 密码加密
func Encode(password, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(password))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// VerifyPassword 密码校验
func VerifyPassword(plaintext, salt, encrypt string) bool {
	return Encode(plaintext, salt) == encrypt
}
