package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// Md5
// md5加密
func Md5(str string) string {
	data := []byte(str)
	rest := fmt.Sprintf("%x", md5.Sum(data))

	return rest
}

// md5 fun v2
func md5V2(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// CryptUserPassword
// 加密用户密码, md5加盐值
func CryptUserPassword(password string, salt string) string {
	return Md5(password + salt)
}

// Salt
// 获取4位密码的盐值
func Salt() string {
	rand.Seed(time.Now().UnixNano()) // 伪随机种子
	baseStr := "abcdefghigklmnopqistuvwxyzABCDEFGHIGKLMNOPQISTUVWXYZ0123456789"
	saltLen := 4
	salt := make([]byte, saltLen)
	for n := 0; n < saltLen; n++ {
		salt[n] = baseStr[rand.Int31n(int32(len(baseStr)))]
	}

	return string(salt)
}

// VerifyUserPassword
// 验证用户的密码是否正确
func VerifyUserPassword(salt string, oldPsd string, oldPassword string) bool {
	password := CryptUserPassword(oldPsd, salt)
	if password == oldPassword {
		return true
	}

	return false
}
