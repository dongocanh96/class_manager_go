package util

import "time"

func CheckDate(time1, time2 time.Time) bool {
	if time1.Year() == time2.Year() && time1.Month() == time2.Month() && time1.Day() == time2.Day() {
		return true
	}

	return false
}
