package myutils

import (
	"os"
	"path/filepath"
	"runtime"
)

//EnvGet 获得环境变量列表
func EnvGet() []string {
	return os.Environ()
}

//EnvGetOsCategory 获得当前系统类型
func EnvGetOsCategory() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

//EnvGetTmpPath 获得系统的临时目录位置
func EnvGetTmpPath() string {
	osString, _ := EnvGetOsCategory()
	switch osString {
	case "windows":
		return "C:/WINDOWS/Temp/"
	case "darwin":
		return "/tmp/"
	case "linux":
		return "/tmp/"
	default:
		panic("the system " + osString + " is not be supported!")
	}
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
