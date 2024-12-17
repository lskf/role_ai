package tests

import (
	"role_ai/infrastructure/utils"
	"testing"
)

func TestCheckString(t *testing.T) {
	data := []string{
		"test",                 // 长度为4
		"12938siod20did942kdi", // 长度为20
		"12938siod20did/.",     // 长度为16, 包含特殊字符
		"28clw9sid9d2kdi",      // 长度为16
	}
	for _, d := range data {
		t.Run(d, func(t *testing.T) {
			if !utils.CheckString(d) {
				t.Fatal("检查失败: 无效的字符串")
			}
			t.Log("检查成功")
		})
	}
}
