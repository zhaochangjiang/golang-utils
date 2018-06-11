package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//SubString 截取字符串
func SubString(s string, start int, length int) string {
	rs := []rune(s)
	if start < 0 {
		start = len(rs) + start
	}
	to := length + start
	key := string(rs[start:to])
	return key
}

//StringLength 获得字符串长度
func StringLength(s string) int {
	rs := []rune(s)
	return len(rs)
}

//StringMd5EqualPHP md5加密字符串,经过验证 此结果与PHP 结果一致
func StringMd5EqualPHP(s string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(s))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
