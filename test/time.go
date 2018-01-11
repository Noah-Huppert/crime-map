package test

import (
	"time"
)

// NewDate creates a time.Time with only a year, month, and day field. In the
// UTC locality.
func NewDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// NewTime creates a time.Time with the following fields: year, month, day,
// hour, and minute. In the UTC locality.
func NewTime(year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
}
