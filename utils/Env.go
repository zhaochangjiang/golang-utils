package utils

import "os"

//EnvGet 获得环境变量列表
func EnvGet() []string {
	return os.Environ()
}

//EnvGetGoHome 获得当前程序的目录
func EnvGetGoHome() string {
	return os.Getenv("GOPATH")
}

//EnvGetPath 获得当前的环境变量
func EnvGetPath() string {
	return os.Getenv("PATH")
}
