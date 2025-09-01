package utils

import (
	"time"
)

func ParseMonthYearToTime(s string) (time.Time, error) {
	return time.Parse("01-2006", s)
}
func TimeToMonthYear(t time.Time) string {
	return t.Format("01-2006")
}
