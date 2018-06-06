package src

//SubString 截取字符串
func SubString(s string, start int, to int) string {
	rs := []rune(s)
	key := string(rs[start:to])
	return key
}
