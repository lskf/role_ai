package utils

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
	"unicode"
)

// TrimString 清除字符串中的空格、换行符、回车符、制表符
func TrimString(str string) string {
	str = strings.TrimSpace(str)
	str = strings.Trim(str, "\n")
	str = strings.Trim(str, "\r")
	str = strings.Trim(str, "\t")
	return str
}

// GenerateNonce 生成随机字符串
func GenerateNonce() (string, error) {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// IsDecimal 判断是否是数字
func IsDecimal(str string) bool {
	// 正则表达式模式
	pattern := `^\d+(\.\d+)?$`

	// 编译正则表达式
	regex := regexp.MustCompile(pattern)

	// 匹配字符串
	return regex.MatchString(str)
}

// CheckString 检查字符串是否同时包含字母和数字，并且长度在8到16个字符之间
func CheckString(s string) bool {
	// 检查长度
	if len(s) < 8 || len(s) > 16 {
		return false
	}

	// 检查是否至少包含一个字母和一个数字
	hasDigit := false
	hasLetter := false
	for _, r := range s {
		if unicode.IsDigit(r) {
			hasDigit = true
		} else if unicode.IsLetter(r) {
			hasLetter = true
		}
		if hasDigit && hasLetter {
			break
		}
	}

	return hasDigit && hasLetter
}
