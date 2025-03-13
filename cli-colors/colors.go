package clicolors

import (
	"time"
)

const (
	Reset     = "\033[0m"  // Reset the color
	Red       = "\033[31m" // Red color
	Green     = "\033[32m" // Green color
	Yellow    = "\033[33m" // Yellow color
	Blue      = "\033[34m" // Blue color
	Magenta   = "\033[35m" // Magenta color
	Cyan      = "\033[36m" // Cyan color
	White     = "\033[37m" // White color
	Bold      = "\033[1m"  // Bold text
	Underline = "\033[4m"  // Underlined text
)

func GetMinutesDifference(t time.Time) float64 {
	now := time.Now()
	diff := t.Sub(now).Abs()
	return diff.Minutes()
}

func GetUnreadColor(c int) string {
	if c > 15 {
		return "\033[38;5;196m"
	} else if c > 5 {
		return "\033[38;5;39m"
	} else if c > 0 {
		return "\033[38;5;82m"
	} else {
		return "\033[38;5;15m"
	}
}

func GetEventColor(t time.Time) string {
	minutes := GetMinutesDifference(t)

	if minutes < 60 {
		return "\033[38;5;39m"
	} else if minutes < 30 {
		return "\033[38;5;82m"
	} else if minutes < 15 {
		return Red
	} else if minutes <= 5 {
		return "\033[38;5;196m"
	} else {
		return "\033[38;5;15m"
	}
}
