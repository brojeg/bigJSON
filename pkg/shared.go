package pkg

import (
	"strconv"
	"strings"
)

func ParseHour(timeString string) int {
	hour, _ := strconv.Atoi(timeString[:len(timeString)-2])
	if strings.HasSuffix(timeString, "PM") && hour != 12 {
		hour += 12
	}
	if strings.HasSuffix(timeString, "AM") && hour == 12 {
		hour = 0
	}
	return hour
}
