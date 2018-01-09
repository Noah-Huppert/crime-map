package date

import (
	"fmt"
	"regexp"
	"time"
)

var pqDateRangeRegexp *regexp.Regexp = regexp.MustCompile("\\(\"(.*)\",\"(.*)\"\\)")

// NewRangeFromStr parses a PostgreSQL date range into 2 time.Time structs. The
// Postgres date range should be passed as a string. An error is returned if
// one occurs. Nil on success.
func NewRangeFromStr(field string) (*time.Time, *time.Time, error) {
	// Extract both dates from field
	matches := pqDateRangeRegexp.FindStringSubmatch(field)

	// If no field is not in postgres date range format
	if matches == nil {
		return nil, nil, fmt.Errorf("field is not in postgres date "+
			"range format, field: %s", field)
	}

	startStr := matches[1]
	endStr := matches[2]

	// Parse strings into dates
	startDate, err := NewTimeFromStr(startStr)
	if err != nil {
		return nil, nil, fmt.Errorf("error converting range start "+
			"date into time.Time: %s", err.Error())
	}

	endDate, err := NewTimeFromStr(endStr)
	if err != nil {
		return nil, nil, fmt.Errorf("error converting range end "+
			"date into time.Time: %s", err.Error())
	}

	// Success
	return startDate, endDate, nil
}

// NewTimeFromStr converts the provided RFC3339 date string into a time.Time. An
// error is returned if one occurs.
func NewTimeFromStr(str string) (*time.Time, error) {
	// Attempt to parse
	t, err := time.Parse("2006-01-02 15:04:05-07", str)
	if err != nil {
		return nil, fmt.Errorf("error converting string into RFC3339 time: %s",
			err.Error())
	}

	// Success
	return &t, nil
}
