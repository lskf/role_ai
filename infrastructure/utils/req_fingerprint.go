package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

// GenerateFingerprint 创建请求指纹
func GenerateFingerprint(r *http.Request) (string, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	bodyStr := string(body)
	hasher := sha256.New()
	if _, err := io.WriteString(hasher, r.Method+r.URL.String()+TrimString(bodyStr)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
