package myutils

import (
	"log"
	"os"
)

//ReName 重命名文件
//file 原文件名
//dest 目标文件名
func ReName(file string, dest string) bool {
	err := os.Rename(file, dest) //重命名 C:\\log\\2013.log 文件为install.txt
	if err != nil {
		//如果重命名文件失败,则输出错误 file rename Error!
		log.Println("file rename Error!" + err.Error())
		//打印错误详细信息
		return false
	}
	return true
}
