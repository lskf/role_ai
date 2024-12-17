package tests

import (
	"role_ai/infrastructure/utils"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := utils.GenerateSalt()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("salt: %s", salt)
}

func TestPwdEncode(t *testing.T) {
	data := []struct {
		password string
	}{
		{"iTm3jhPp51Am81deK3Gs"},
	}

	for _, d := range data {
		t.Run(d.password, func(t *testing.T) {
			salt, err := utils.GenerateSalt()
			if err != nil {
				t.Fatal(err.Error())
			}
			encrypt := utils.Encode(d.password, salt)
			t.Logf("password: %s, salt: %s, encrypt: %s", d.password, salt, encrypt)
		})
	}
}

func TestPwdVerify(t *testing.T) {
	data := []struct {
		plaintext string
		salt      string
		encode    string
	}{
		{
			plaintext: "123456",
			salt:      "FkFw7gVsowiCNCFesAiMC8BQtVPfqH8FmjpDjrxq1W9F",
			encode:    "eBucOtukB1eJuE9C3LyQYqHImHFZ1/LyNp7P0pvsj6w=",
		},
	}

	for _, d := range data {
		t.Run(d.plaintext, func(t *testing.T) {
			if !utils.VerifyPassword(d.plaintext, d.salt, d.encode) {
				t.Fatalf("verify failed")
			}
			t.Logf("verify success")
		})
	}
}
