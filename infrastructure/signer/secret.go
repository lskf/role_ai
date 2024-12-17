package signer

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"role_ai/infrastructure/utils"
	"time"
)

func GenerateAppId() string {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s", time.Now().Format("06"), new(big.Int).SetBytes(bytes).String())
}

func GenerateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "TO-" + string(utils.Base58Encode(bytes)), nil
}
