package myutils

import (
	"os"
	"path/filepath"
)

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

//EnvGetCurrentPath 获得当前的目录
func EnvGetCurrentPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir, err
}
