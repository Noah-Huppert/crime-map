package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Noah-Huppert/crime-map/crime"
)

// headerDateRangeExpr is the regexp used to match a report header's first line
var headerDateRangeExpr *regexp.Regexp = regexp.MustCompile("^From [A-Z][a-z]+ [0-9]{1,2}, [0-9]{4} to [A-Z][a-z]+ [0-9]{1,2}, [0-9]{4}\\.$")

// dateExpr is the regexp used to match a date in the pdf report
var dateExpr *regexp.Regexp = regexp.MustCompile("^([0-9]{2})\\/([0-9]{2})\\/([0-9]{2}) - [A-Z]+ at ([0-9]{2}):([0-9]{2})$")

// dateRangeExpr is the regexp used to extract 2 dates in a report date range
var dateRangeExpr *regexp.Regexp = regexp.MustCompile("^(.*[0-9]) - ([0-9].*)$")

// footerPageNumExpr holds the regexp used to match page numbers at the bottom
// of the header
var footerPageNumExpr *regexp.Regexp = regexp.MustCompile("^[0-9]+$")

// Field labels
const fieldLabelReported string = "Date Reported:"
const fieldLabelIncidents string = "Incident(s):"
const fieldLabelOccurred string = "Date and Time Occurred From - Occurred To:"
const fieldLabelFix string = "Disposition:"
const fieldLabelCrimeCount string = "Incident(s) Listed."

// DrexelParser implements the Parser interface for Drexel University Clery
// crime logs
type DrexelParser struct{}

// NewDrexelParser creates a new DrexelParser instance
func NewDrexelParser() *DrexelParser {
	return &DrexelParser{}
}

// Parse interprets a pdf's text fields into Crime structs. For the style of
// report Drexel University releases.
func (p DrexelParser) Parse(fields []string) ([]crime.Crime, error) {
	// crimes holds all parsed Crimes
	crimes := []crime.Crime{}

	// c holds the Crime struct currently being parsed
	var c crime.Crime

	// skip indicates how many fields the parser should skip
	var skip int

	// consume indicates how many fields the parser should consume, if
	// multiple fields need to be consumed in a row
	var consume int

	// consumeGlob1 indicates if the date reported, location, report ID,
	// and incidents field values come after the skipping is done
	var consumeGlob1 bool

	// consumeDOccurred indicates if the date occurred field follows the skipped
	// fields
	var consumeDOccurred bool

	// consumeDesc indicates if the synopsis field is being consumed
	var consumeDesc bool

	// consumeFix indicates that the remediation field is coming up next
	var consumeFix bool

	// consumeCrimeCount indicates that the total crime count field is
	// coming up next
	var consumeCrimeCount bool

	// Loop through fields
	for _, field := range fields {
		// Check if we are skipping fields
		if skip > 0 {
			skip -= 1
			continue
		}

		// Check if first line of header
		if headerDateRangeExpr.MatchString(field) {
			skip = 3
		} else if footerPageNumExpr.MatchString(field) { // Check if
			// first line of footer
			skip = 5
		} else if consumeGlob1 { // Check if we are consuming glob 1
			// If consuming date reported field
			if consume == 4 {
				d, err := parseDate(field)
				if err != nil {
					return crimes, fmt.Errorf("error parsing"+
						" reported at field: %s",
						err.Error())
				}

				c.DateReported = *d
				consume--
			} else if consume == 3 { // If consuming location field
				c.Location = field
				consume--
			} else if consume == 2 { // If consuming report ID field
				// Split by dash
				parts := strings.Split(field, "-")

				// Check correct number of parts
				if len(parts) != 2 {
					return crimes, fmt.Errorf("report ID "+
						"field has incorrect number of"+
						" parts, field: %s, parts: %d, "+
						" expected parts: 2",
						field, len(parts))
				}

				// Parse both ids
				id, err := strconv.ParseUint(parts[0], 10, 64)
				if err != nil {
					return crimes, fmt.Errorf("error parsing "+
						"report super ID into uint: %s",
						err.Error())
				}
				c.ReportSuperID = uint(id)

				id, err = strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return crimes, fmt.Errorf("error parsing "+
						"report Id into uint: %s",
						err.Error())
				}
				c.ReportID = uint(id)

				consume--
			} else if consume == 1 { // If consuming incidents field
				c.Incidents = []string{field}

				// And reset glob 1 parser flags
				consumeGlob1 = false
				consume = 0
			}
		} else if consumeDOccurred { // Check if we are consuming date
			// occurred

			// Split dates
			matches := dateRangeExpr.FindStringSubmatch(field)

			// Check correct number of dates
			if len(matches) != 3 {
				return crimes, fmt.Errorf("error parsing date "+
					"occurred, incorrect number of dates, "+
					"field: %s, expected 2, got: %d",
					field, len(matches)-1)
			}

			// Parse dates
			start, err := parseDate(matches[1])
			if err != nil {
				return crimes, fmt.Errorf("error parsing occurred"+
					" start date, field: %s, err: %s",
					field, err.Error())
			}

			end, err := parseDate(matches[2])
			if err != nil {
				return crimes, fmt.Errorf("error parsing occurred"+
					" end date, field: %s, err: %s",
					field, err.Error())
			}

			// Save
			c.DateOccurredStart = *start
			c.DateOccurredEnd = *end

			consumeDOccurred = false
		} else if consumeDesc { // Check if consuming synopsis
			// Check if end of consuming synopsis
			if field == fieldLabelFix {
				// End of synopsis field values
				consumeDesc = false

				// Beginning of fix field
				consumeFix = true
			} else { // Otherwise, consume
				c.Descriptions = append(c.Descriptions, field)
			}
		} else if consumeFix { // Check if we should consume the
			// remediation field
			c.Remediation = field
			consumeFix = false

			// And add crime to list
			crimes = append(crimes, c)
			c = crime.Crime{}
		} else if field == fieldLabelReported { // Check if beginning
			// of glob 1
			skip = 2
			consumeGlob1 = true
			consume = 4
		} else if field == fieldLabelIncidents { // Check if date
			// occurred field will come next
			consumeDOccurred = true
		} else if field == fieldLabelOccurred { // Check if synopsis
			// field will come next
			skip = 1
			consumeDesc = true
			c.Descriptions = []string{}
		} else if field == fieldLabelCrimeCount { // Check if incidents
			// listed count comes next
			consumeCrimeCount = true
		} else if consumeCrimeCount { // Check if consuming crime count
			count, err := strconv.Atoi(field[1:])
			if err != nil {
				return crimes, fmt.Errorf("error parsing number"+
					" of listed crimes: %s", err.Error())
			}

			consumeCrimeCount = false

			// Check count matches
			if len(crimes) != count {
				return crimes, fmt.Errorf("number of listed "+
					"crimes and number of crimes parsed "+
					"does not match: listed: %d, parsed: %d",
					count,
					len(crimes))
			}
		} else {
			fmt.Printf("UNKNOWN>>%s\n", field)
		}
	}

	return crimes, nil

}

// parseDate Creates a time struct from a drexel date on a report. The offset
// var specifies the offset to add to match indexes. An error is returned if one
// occurs, nil otherwise.
func parseDate(field string) (*time.Time, error) {
	matches := dateExpr.FindStringSubmatch(field)

	year, err := strconv.ParseInt(matches[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing date year: %s",
			err.Error())
	}

	month, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing date month: %s",
			err.Error())
	}

	day, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing date day: %s",
			err.Error())
	}

	hour, err := strconv.ParseInt(matches[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing date hour: %s",
			err.Error())
	}

	minute, err := strconv.ParseInt(matches[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing date minute: %s",
			err.Error())
	}

	d := time.Date(int(year),
		time.Month(month),
		int(day),
		int(hour),
		int(minute),
		0, 0, time.UTC)

	return &d, nil
}
