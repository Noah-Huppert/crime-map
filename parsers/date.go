package parsers

import (
	"fmt"
	"regexp"

	"github.com/Noah-Huppert/crime-map/date"
	"github.com/Noah-Huppert/crime-map/models"
)

// DateRangeParser parses date range fields in the form:
//
//	From <month abbrev> <date>, <year> to <month abbrev> <date>, <year>.
//
// Where month abbreviations are the first 3 letters of a month's name, with a
// capital first letter.
//
// Dates can be either 1 or 2 digits. Years must be in their full 4 digit form.
type DateRangeParser struct{}

// rangeExpr is a regexp used to match a date range found in a crime report
var rangeExpr *regexp.Regexp = regexp.MustCompile("^From ([A-Z][a-z]+) ([0-9]{1,2}), ([0-9]{4}) to ([A-Z][a-z]+) ([0-9]{1,2}), ([0-9]{4})\\.$")

// Parse implements the Parser.Parse method for the DateRangeParser
func (p DateRangeParser) Parse(i uint, fields []string, report *models.Report,
	crime *models.Crime) (uint, error) {

	if matches := rangeExpr.FindStringSubmatch(fields[i]); matches != nil {
		// Convert start date
		startRange, err := date.ParseDate(matches[3:6])
		if err != nil {
			return 0, fmt.Errorf("error converting start "+
				"header date to time.Time: %s",
				err.Error())
		}
		report.RangeStartDate = startRange

		// Convert end date
		endRange, err := date.ParseDate(matches[3:6])
		if err != nil {
			return 0, fmt.Errorf("error converting end "+
				"header date to time.Time: %s",
				err.Error())
		}
		report.RangeEndDate = endRange

		// Mark as parsed and skip fields
		return 3, nil

	}

	return 0, nil
}

func (p DateRangeParser) String() string {
	return "DateRangeParser"
}
