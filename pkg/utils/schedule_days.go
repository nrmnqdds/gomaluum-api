package utils

import (
	"strings"
)

func GetScheduleDays(day string) uint8 {
	if strings.Contains(day, "SUN") {
		return 0
	} else if strings.Contains(day, "MON") || day == "M" {
		return 1
	} else if strings.Contains(day, "TUE") || day == "T" {
		return 2
	} else if strings.Contains(day, "WED") || day == "W" {
		return 3
	} else if strings.Contains(day, "THUR") || day == "TH" {
		return 4
	} else if strings.Contains(day, "FRI") || day == "F" {
		return 5
	} else if strings.Contains(day, "SAT") {
		return 6
	}
	return 7
}
