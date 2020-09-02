package utility

import (
	"time"
)

// FormatDate formats a date.
func FormatDate(date time.Time) string {
	return date.Format("Mon _2 Jan 2006 15:04:05")
}
