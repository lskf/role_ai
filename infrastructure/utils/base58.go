package utils

import (
	"bytes"
	"math/big"
)

// Base58 字符集
var base58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// Base58Encode 编码字节数组到 base58
func Base58Encode(input []byte) []byte {
	var result []byte

	// 将 input 转换为一个大整数
	x := big.NewInt(0).SetBytes(input)

	// base 58
	base := big.NewInt(58)

	// 循环直到整数为 0
	zero := big.NewInt(0)
	mod := &big.Int{}
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod) // 对 base 取模
		result = append(result, base58Alphabet[mod.Int64()])
	}

	// 添加前导零
	for _, b := range input {
		if b != 0 {
			break
		}
		result = append(result, base58Alphabet[0])
	}

	// 反转切片，因为 DivMod 是从低位开始计算的
	reverse(result)
	return result
}

// Base58Decode 解码 base58 字符串到字节数组
func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	for _, b := range input {
		value := bytes.IndexByte(base58Alphabet, b)
		result.Mul(result, big.NewInt(58))
		result.Add(result, big.NewInt(int64(value)))
	}

	// 解码前导零
	var decoded []byte
	for _, b := range input {
		if b != base58Alphabet[0] {
			break
		}
		decoded = append(decoded, 0)
	}

	decoded = append(decoded, result.Bytes()...)
	return decoded
}

// reverse 反转字节切片
func reverse(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
