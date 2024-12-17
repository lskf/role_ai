package tools

import (
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomIntInRange 生成指定范围的随机Int
func RandomIntInRange(min, max int) int {
	return random.Intn(max-min+1) + min
}

// RandomInt64InRange 生成指定范围的随机Int64
func RandomInt64InRange(min, max int64) int64 {
	return random.Int63n(max-min+1) + min
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}
