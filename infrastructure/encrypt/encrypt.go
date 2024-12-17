package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// GenerateEncryptKey 生成一个随机的 AES 密钥
func GenerateEncryptKey() (string, error) {
	key := make([]byte, 32) // AES-256，所以使用 32 字节长度的密钥
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// Encrypt 对文本进行 AES 加密, 加密的方式为 AES-256-CBC
func Encrypt(plaintext, key string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 原始文本转为字节
	byteText := []byte(plaintext)

	// 块加密的数据必须对齐块大小
	blockSize := block.BlockSize()
	byteText = PKCS7Padding(byteText, blockSize)

	// 初始化向量 IV 必须是唯一的，但不需要保密
	ciphertext := make([]byte, blockSize+len(byteText))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 使用 CBC 模式加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], byteText)

	// 将二进制密文转化为 base64
	// 返回结果是 IV 和密文的组合，这是常见的做法
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 对 AES 加密的文本进行解密
func Decrypt(cryptoText, key string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// IV 从加密的数据中提取出来
	blockSize := block.BlockSize()
	if len(ciphertext) < blockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:blockSize]
	ciphertext = ciphertext[blockSize:]

	// CBC 模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 移除填充
	ciphertext = PKCS7UnPadding(ciphertext)
	return string(ciphertext), nil
}

// PKCS7Padding 补全块
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 删除补全块
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
