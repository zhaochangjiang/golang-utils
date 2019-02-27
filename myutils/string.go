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
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
)

//StringPregReplace 正则表达式替换
func StringPregReplace(preg, replaceToString, fromString string) string {

	return fromString
}

//StringSub 截取字符串
func StringSub(s string, start ...int) string {
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

//StringImplode 连接字符串
func StringImplode(separator string, array *[]string) string {
	var i = 0
	res := ""
	for _, v := range *array {
		if i == 0 {
			res += v
		} else {
			res += separator + v
		}
		i++
	}
	return res
}

//StringGUID 生成Guid字串
func StringGUID() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return StringMd5EqualPHP(base64.URLEncoding.EncodeToString(b))
}

//StringReplace 替换字符串
func StringReplace(source string, dividString *map[string]string) string {

	sliceRes := make([]stringReplaceStruct, 0)
	obj := newStringReplaceStruct()
	obj.replaceString = ""
	obj.nowString = source
	sliceRes = append(sliceRes, *obj)

	for k, v := range *dividString {
		sliceRes = *(everyStringDivid(k, v, &sliceRes))
	}
	return concatString(&sliceRes)
}

//StringHtmlspecialchars  转换字符串
func StringHtmlspecialchars(s string) string {
	return StringReplace(s, &map[string]string{
		"&": "&#38;", "\"": "&#34;", "<": "&#60;", ">": "&#62;", "'": "&#39;",
	})
}

/*******************private method and struct**********************/
func everyStringDivid(dividString, replace string, sliceResPoniter *[]stringReplaceStruct) *[]stringReplaceStruct {
	sliceRes := make([]stringReplaceStruct, 0)
	for _, v := range *sliceResPoniter {
		ls := strings.Split(v.nowString, dividString)

		convertStringToStringReplaceStruct(&sliceRes, &ls, replace, v.replaceString)
	}

	return &sliceRes
}

func concatString(srs *[]stringReplaceStruct) string {
	s := ""
	for k, v := range *srs {
		if k == 0 {
			s += v.nowString
		} else {
			s += v.replaceString + v.nowString
		}
	}
	return s
}

func convertStringToStringReplaceStruct(srs *[]stringReplaceStruct, sa *[]string, replaceString, replace string) {

	for k, v := range *sa {
		obj := newStringReplaceStruct()
		//如果是开始的时候才需要
		if k != 0 {
			obj.replaceString = replaceString // + replace
		} else {
			obj.replaceString = replace
		}

		obj.nowString = v
		*srs = append(*srs, *obj)
	}
}

type stringReplaceStruct struct {
	replaceString string
	nowString     string
}

func newStringReplaceStruct() *stringReplaceStruct {
	obj := &stringReplaceStruct{}
	return obj
}
