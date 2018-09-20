package lib

import "fmt"

//Dsc 一个减法实现
// 返回a-b的值
func Dsc(b int, c int) int {
	//Output:1+2=3
	return b - c
}

func Example() {
	sum := Add(1, 2)
	fmt.Println("1+2=", sum)
	//Output:
	//1+2=3
}
