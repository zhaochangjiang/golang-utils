package myutils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
)

//SubString 截取字符串
func SubString(s string, start ...int) string {
	cap := cap(start)
	if cap < 1 {
		panic("the params start must be 1 or 2 parameters!")
	}
	rs := []rune(s)
	len := len(rs)
	if cap < 2 {
		start = append(start, len-start[0])
	}
	if start[0] < 0 {
		start[0] = len + start[0]
	}
	to := start[1] + start[0]
	key := string(rs[start[0]:to])
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

//StringBase64Encode base64加密字符
func StringBase64Encode(raw []byte) []byte {
	var encoded bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &encoded)
	encoder.Write(raw)
	encoder.Close()
	return encoded.Bytes()
}

//HMAC_SHA1加密
func StringHMACSHA1(s []byte, secret string) []byte {
	//sha1//	YourSecretKey := "Wgyokg4EC2dJKwOYwBu3zrZHgfnOYHWi"
	//hmac ,use sha1
	key := []byte(secret)
	mac := hmac.New(sha1.New, key)
	mac.Write(s)
	res := mac.Sum(nil)
	return res
}

type Charset string

//StringUTF8EncodingOf
func StringUTF8EncodingOf(byteContent []byte, charset Charset) []byte {
	const (
		UTF8    = Charset("UTF-8")
		GB18030 = Charset("GB18030")
	)
	var str []byte
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byteContent)
		str = decodeBytes
	case UTF8:
		fallthrough
	default:
		str = byteContent
	}

	return str
}

//StringGUID 生成Guid字串
func StringGUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return StringMd5EqualPHP(base64.URLEncoding.EncodeToString(b))
}
