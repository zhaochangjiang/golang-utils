package utils

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
