package utils

import "crypto/rand"

// 生成盐
func GenerateSalt(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	bytes := make([]byte, length+(length/4))
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes[:length] { // 修改此处，仅遍历有效长度范围的字节
		result[i] = charset[b%byte(len(charset))]
	}
	return string(result), nil
}
