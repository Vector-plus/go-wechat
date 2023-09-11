package secure

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// 生成秘钥
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 加密
func MakePassword(pd, salt string) string {
	return Md5Encode(pd + salt)
}

// 解密判断
func ValidPassword(pd, salt, password string) bool {
	return password == MakePassword(pd, salt)
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		index, _ := rand.Int(rand.Reader, big.NewInt(60))
		b[i] = defaultLetters[index.Int64()]
	}
	return string(b)
}
