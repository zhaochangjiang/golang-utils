package main

import (
	"log"

	"github.com/zhaochangjiang/golang-utils/myutils"
)

func main() {
	params := make(map[string]string)
	params["a"] = "2147483646"
	params["b"] = "18608106929"
	rules := []myutils.ValidateRule{
		{Key: "a", Set: true, Empty: true, Alias: "参数a", Category: "uint64", Min: -1, Max: 2147483645},
		// {Key: "b", Set: true, Empty: true, Category: "int", DefaultValue: "adefault", Alias: "参数B", Preg: "1[0-9]{10}"},
	}
	result := myutils.NewValidate(&params, &rules).Run()
	log.Println(result)
	log.Println(params)
}
