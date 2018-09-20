package myutils

import (
	"testing"
)

func Test_StringReplace(t *testing.T) {
	p := make(map[string]string)
	p["123"] = "kkk"
	p["kafs"] = "mmm"

	s := "测试函数123afsdvvvvsdf1231测试dfsd"
	r := StringReplace(s, &p)
	if r != "abc123afsdasdf1231asdfsd" {
		t.Error("test error")
	}
}
