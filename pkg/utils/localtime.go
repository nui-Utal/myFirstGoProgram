package utils

import (
	"time"
)

func GetLocalTime() string {
	currentTime := time.Now()

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return ""
	}

	currentTime = time.Now().In(loc)

	return currentTime.Format("2006-01-02 15:04:05")
}
