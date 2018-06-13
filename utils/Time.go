package utils

import (
	"time"
)

//GetTimeNow
func TimeGetNowStruct() time.Time {
	return time.Now()
}

//TimeNow 获得当前的时间 YYYY-MM-dd HH:ii:ss
func TimeNow() string {
	return TimeGetNowStruct().Format("2006-01-02 15:04:05")
}

//TimeGetUnixNano 获得当前的unix时间戳
func TimeGetUnixNano() int64 {
	return TimeGetNowStruct().Unix()
}
