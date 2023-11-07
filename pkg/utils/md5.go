package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encrypt(str string) string {
	hash := md5.Sum([]byte(str))
	encryptedStr := hex.EncodeToString(hash[:])
	return encryptedStr
}
