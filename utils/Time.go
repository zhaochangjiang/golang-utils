package utils

import "time"

//TimeNow 获得当前的时间 YYYY-MM-dd HH:ii:ss
func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
