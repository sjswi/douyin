package tools

import (
	"fmt"
	"time"
)

func GetMonthAndDay(now time.Time) string {
	month := now.Month()
	day := now.Day()
	return fmt.Sprintf("%d-%d", month, day)
}
