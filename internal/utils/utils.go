package utils

import "time"

func SleepSeconds(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
