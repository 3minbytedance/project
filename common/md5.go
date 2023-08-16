package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// Md5Encoder md5加密
func Md5Encoder(data string) string {
	md := md5.New()
	md.Write([]byte(data))
	tempStr := md.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// MakePassword 加密密码
func MakePassword(originalPwd string, randNum string) string {
	return Md5Encoder(originalPwd + randNum)
}

// CheckPassword 判断密码是否相等
func CheckPassword(originalPwd string, randNum string, password string) bool {
	return password == MakePassword(originalPwd, randNum)
}

func MakeToken() string {
	timeStr := fmt.Sprintf("%d", time.Now().Unix())
	return Md5Encoder(timeStr)
}
