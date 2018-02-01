package date

import (
	"fmt"
	"strconv"
	"time"
)

// ParseDateParts extracts date information from a provided slice of strings.
// The slice's first element should be the month, followed by the day and then
// the year. The month and day can be either 1 or 2 digits. However the year
// must be in full 4 digit form.
//
// The parsed time will be returned. Along with an error if one occurs, nil on
// success.
func ParseDate(parts []string) (*time.Time, error) {
	// Check parts length
	if len(parts) != 3 {
		return nil, fmt.Errorf("provided parts slice should contain "+
			"3 elements, had: %d", len(parts))
	}

	// Extract info
	monthStr := parts[0]
	dayStr := parts[1]
	yearStr := parts[2]

	// Parse month
	month, err := ParseMonthAbbrv(monthStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing month "+
			"abbreviation into month number: %s", err.Error())
	}

	// Parse day
	day, err := strconv.ParseUint(dayStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing day: %s",
			err.Error())
	}

	// Parse year
	year, err := strconv.ParseUint(yearStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing year: %s",
			err.Error())
	}

	// Success
	t := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0,
		time.UTC)

	return &t, nil

}

// monthAbbrvs holds valid month abbreviations
var monthAbbrvs []string = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// ParseMonthAbbrv extracts and returns the month number from a 3 character
// abbreviation. An error is returned if one occurs, nil on success.
func ParseMonthAbbrv(abbrv string) (uint, error) {
	// Linear search for valid abbreviation
	for i, val := range monthAbbrvs {
		if val == abbrv {
			return uint(i + 1), nil
		}
	}

	// If none found
	return 0, fmt.Errorf("error parsing month abbreviation, unknown "+
		"value: %s", abbrv)
}
