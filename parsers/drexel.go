package parsers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Noah-Huppert/crime-map/geo"
	"github.com/Noah-Huppert/crime-map/models"
)

// DrexelUName holds the Drexel University's name
const DrexelUName string = "Drexel University"

// headerDateRangeExpr is the regexp used to match a report header's first line
var headerDateRangeExpr *regexp.Regexp = regexp.MustCompile("^From ([A-Z][a-z]+) ([0-9]{1,2}), ([0-9]{4}) to ([A-Z][a-z]+) ([0-9]{1,2}), ([0-9]{4})\\.$")

// dateExpr is the regexp used to match a date in the pdf report
var dateExpr *regexp.Regexp = regexp.MustCompile("^([0-9]{2})\\/([0-9]{2})\\/([0-9]{2}) - [A-Z]+ at ([0-9]{2}):([0-9]{2})$")

// dateRangeExpr is the regexp used to extract 2 dates in a report date range
var dateRangeExpr *regexp.Regexp = regexp.MustCompile("^(.*[0-9]) - ([0-9].*)$")

// footerPageNumExpr holds the regexp used to match page numbers at the bottom
// of the header
var footerPageNumExpr *regexp.Regexp = regexp.MustCompile("^[0-9]+$")

// monthAbbrvs holds valid month abbreviations
var monthAbbrvs []string = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun",
	"Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

// Field labels
const fieldLabelReported string = "Date Reported:"
const fieldLabelIncidents string = "Incident(s):"
const fieldLabelOccurred string = "Date and Time Occurred From - Occurred To:"
const fieldLabelFix string = "Disposition:"
const fieldLabelCrimeCount string = "Incident(s) Listed."

// DrexelParser implements the Parser interface for Drexel University Clery
// crime logs
type DrexelParser struct {
	// geoCache is used to cache GeoLoc queryies to the database
	geoCache *geo.GeoCache

	// parsed indicates if a report has been parsed yet
	parsed bool

	// crimes holds the Crimes which were parsed from a report, empty if
	// parsed == false
	crimes []models.Crime

	// startRange holds the start of the time range which the report covers
	startRange time.Time

	// endRange holds the end of the time range which the report covers
	endRange time.Time
}

// NewDrexelParser creates a new DrexelParser instance
func NewDrexelParser(geoCache *geo.GeoCache) *DrexelParser {
	return &DrexelParser{
		geoCache: geoCache,
		parsed:   false,
		crimes:   []models.Crime{},
	}
}

// Range implements the Range method for Parser
func (p DrexelParser) Range() (*time.Time, *time.Time, error) {
	// If not parsed
	if !p.parsed {
		return nil, nil, ErrReportNotParsed
	}

	// If parsed
	return &p.startRange, &p.endRange, nil
}

// Parse interprets a pdf's text fields into Crime structs. For the style of
// report Drexel University releases.
func (p *DrexelParser) Parse(fields []string) ([]models.Crime, error) {
	// Check if already parsed
	if p.parsed {
		// Return results
		return p.crimes, ErrReportParsed
	}

	// c holds the Crime struct currently being parsed
	var c models.Crime

	// skip indicates how many fields the parser should skip
	var skip int

	// consume indicates how many fields the parser should consume, if
	// multiple fields need to be consumed in a row
	var consume int

	// haveConsumedHeaderDate indicates that the header date range has
	// already been parsed
	var haveConsumedHeaderDate bool

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

	// pageNum holds the current page number
	pageNum := 0

	// Loop through fields
	for _, field := range fields {
		// Check if we are skipping fields
		if skip > 0 {
			skip -= 1
			continue
		}

		// Check if first line of header
		if matches := headerDateRangeExpr.FindStringSubmatch(field); matches != nil {
			// Header
			pageNum += 1
			skip = 3

			// If already parsed
			if haveConsumedHeaderDate {
				continue
			}

			// Convert start date
			startMonthStr := matches[1]
			startDateStr := matches[2]
			startYearStr := matches[3]

			startRange, err := p.parseHeaderDate(startMonthStr,
				startDateStr, startYearStr)
			if err != nil {
				return p.crimes, fmt.Errorf("error converting "+
					"start header date to time.Time: %s",
					err.Error())
			}
			p.startRange = *startRange

			// Convert end date
			endMonthStr := matches[4]
			endDateStr := matches[5]
			endYearStr := matches[6]

			endRange, err := p.parseHeaderDate(endMonthStr,
				endDateStr, endYearStr)
			if err != nil {
				return p.crimes, fmt.Errorf("error converting "+
					"end header date to time.Time: %s",
					err.Error())
			}
			p.endRange = *endRange

			// Mark as parsed and skip fields
			haveConsumedHeaderDate = true
		} else if footerPageNumExpr.MatchString(field) { // Check if
			// first line of footer
			skip = 5
		} else if consumeGlob1 { // Check if we are consuming glob 1
			// If consuming date reported field
			if consume == 4 {
				d, err := parseDate(field)
				if err != nil {
					return p.crimes, fmt.Errorf("error parsing"+
						" reported at field: %s",
						err.Error())
				}

				c.Page = pageNum
				c.DateReported = *d
				consume--
			} else if consume == 3 { // If consuming location field
				// Gets GeoLoc with just a populated ID field.
				// This allows us to set the crime foreign key,
				// but not know anything about the location
				loc, err := p.geoCache.InsertIfNew(field)

				if err != nil {
					return p.crimes, fmt.Errorf("error "+
						"getting cached GeoLoc: %s",
						err.Error())
				}

				// Set GeoLoc FK
				c.GeoLocID = loc.ID

				consume--
			} else if consume == 2 { // If consuming report ID field
				// Split by dash
				parts := strings.Split(field, "-")

				// Check correct number of parts
				if len(parts) != 2 {
					return p.crimes, fmt.Errorf("report ID "+
						"field has incorrect number of"+
						" parts, field: %s, parts: %d, "+
						" expected parts: 2",
						field, len(parts))
				}

				// Parse both ids
				id, err := strconv.ParseUint(parts[0], 10, 64)
				if err != nil {
					return p.crimes, fmt.Errorf("error parsing "+
						"report super ID into uint: %s",
						err.Error())
				}
				c.ReportSuperID = uint(id)

				id, err = strconv.ParseUint(parts[1], 10, 64)
				if err != nil {
					return p.crimes, fmt.Errorf("error parsing "+
						"report Id into uint: %s",
						err.Error())
				}
				c.ReportSubID = uint(id)

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
				return p.crimes, fmt.Errorf("error parsing date "+
					"occurred, incorrect number of dates, "+
					"field: %s, expected 2, got: %d",
					field, len(matches)-1)
			}

			// Parse dates
			start, err := parseDate(matches[1])
			if err != nil {
				return p.crimes, fmt.Errorf("error parsing occurred"+
					" start date, field: %s, err: %s",
					field, err.Error())
			}

			end, err := parseDate(matches[2])
			if err != nil {
				return p.crimes, fmt.Errorf("error parsing occurred"+
					" end date, field: %s, err: %s",
					field, err.Error())
			}

			// Check start date is after end date
			if start.After(*end) {
				// If so, add 12 hours to end date
				fixedEnd := end.Add(time.Hour *
					time.Duration(12))

				// Check again
				if start.After(fixedEnd) {
					// We don't know how to fix, error
					return p.crimes, fmt.Errorf("error "+
						"parsing occurred date, start "+
						"date is before end date, "+
						"after correction, field: %s",
						field)
				}

				// Note parse error
				pErr := models.ParseError{
					Field:    "date_occurred",
					Original: field,
					Corrected: fmt.Sprintf("%s - %s",
						start.String(),
						fixedEnd.String()),
					ErrType: models.TypeBadRangeEnd,
				}

				// Save parse error
				c.ParseErrors = append(c.ParseErrors, pErr)

				// If success, replace
				end = &fixedEnd
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
			p.crimes = append(p.crimes, c)
			c = models.Crime{}
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
				return p.crimes, fmt.Errorf("error parsing number"+
					" of listed crimes: %s", err.Error())
			}

			consumeCrimeCount = false

			// Check count matches
			if len(p.crimes) != count {
				return p.crimes, fmt.Errorf("number of listed "+
					"crimes and number of crimes parsed "+
					"does not match: listed: %d, parsed: %d",
					count,
					len(p.crimes))
			}
		} else {
			return p.crimes, fmt.Errorf("error parsing field, "+
				"unknown value: %s", field)
		}
	}

	// Success
	p.parsed = true
	return p.crimes, nil

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

// parseHeaderDate parses a time.Time from a report header. A Date will be
// returned. Along with an error if one occurs. Nil on success.
func (p DrexelParser) parseHeaderDate(monthStr string, dateStr string, yearStr string) (*time.Time, error) {
	// Parse month
	month, err := p.parseMonthAbbrv(monthStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing month "+
			"abbreviation into month number: %s", err.Error())
	}

	// Parse date
	date, err := strconv.ParseUint(dateStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %s",
			err.Error())
	}

	// Parse year
	year, err := strconv.ParseUint(yearStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parsing year: %s",
			err.Error())
	}

	// Success
	t := time.Date(int(year), time.Month(month), int(date), 0, 0, 0, 0,
		time.UTC)
	return &t, nil
}

// parseMonthAbbrv extracts and returns the month number from a 3 character
// abbreviation. An error is returned if one occurs, nil on success.
func (p DrexelParser) parseMonthAbbrv(abbrv string) (uint, error) {
	// Linear search valid abbreviations
	for i, val := range monthAbbrvs {
		if val == abbrv {
			return uint(i + 1), nil
		}
	}

	// If none found
	return 0, fmt.Errorf("error parsing month abbreviation, unknown value"+
		": %s", abbrv)
}
