package parsers

// NewDrexelRunner creates a new ParserRunner populated with the correct parsers
// needed extract all crime report information.
func NewDrexelRunner() *ParserRunner {
	r := NewParserRunner()
	r.Add(DateRangeParser{})

	return r
}

/*
import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/Noah-Huppert/crime-map/errs"
	"github.com/Noah-Huppert/crime-map/geo"
	"github.com/Noah-Huppert/crime-map/models"
)

// DrexelUName holds the Drexel University's name
const DrexelUName string = "Drexel University"

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

// errNotHeaderDateRange indicates that the provided field was not a date range
// present in the report header
var errNotHeaderDateRange error = errors.New("provided field was not a header" +
	" date range")

// DrexelParser implements the Parser interface for Drexel University Clery
// crime logs
type DrexelParser struct {
	// logger is used to output debug information
	logger *log.Logger

	// geoCache is used to cache GeoLoc queryies to the database
	geoCache *geo.GeoCache

	// parsedCrimes indicates if a report's crime models have been parsed
	// out yet
	parsedCrimes bool

	// parsedRange indicates if a report's date range has been parsed
	// out yet
	parsedRange bool

	// crimes holds the Crimes which were parsed from a report, empty if
	// parsedCrimes == false
	crimes []models.Crime

	// startRange holds the start of the time range which the report covers
	startRange *time.Time

	// endRange holds the end of the time range which the report covers
	endRange *time.Time
}

// NewDrexelParser creates a new DrexelParser instance
func NewDrexelParser(geoCache *geo.GeoCache) *DrexelParser {
	return &DrexelParser{
		logger:       log.New(os.Stdout, "parsers/drexel", 0),
		geoCache:     geoCache,
		parsedCrimes: false,
		parsedRange:  false,
		crimes:       []models.Crime{},
	}
}

// Range implements the Range method for Parser. It parses the fields far
// enough to determine the date range the report covers
func (p *DrexelParser) Range(fields []string) (*time.Time, *time.Time, error) {
	// Check if already parsed range
	if p.parsedRange {
		// If so, return
		return p.startRange, p.endRange, nil
	}

	// Loop through fields until we parse a header date range
	for _, field := range fields {
		// If parsed header date range
		if _, err := p.parseHeaderRange(field); err != errNotHeaderDateRange {
			// If parse error
			if err != nil {
				return nil, nil, fmt.Errorf("error parsing "+
					"header date range: %s", err.Error())
			}

			// Success
			return p.startRange, p.endRange, nil
		}
	}

	// If looped through all fields and not found, error
	return nil, nil, errors.New("error finding header date range, not found")
}

// Count returns the number of crimes parsed in the report.
func (p DrexelParser) Count() (uint, error) {
	// Check if not parsed yet
	if !p.parsedCrimes {
		return 0, errs.ErrNotParsed
	}

	// If parsed, return count
	l := len(p.crimes)

	return uint(l), nil
}

// Parse interprets a pdf's text fields into Crime structs. For the style of
// report Drexel University releases.
func (p *DrexelParser) Parse(i uint, fields []string, report *models.Report, crime *models.Crime) (uint, error) {
		// Check if already parsed
		if p.parsedCrimes {
			// Return results
			return p.crimes, errs.ErrParsed
		}

		// c holds the Crime struct currently being parsed
		var c models.Crime

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
			if newSkip, err := p.parseHeaderRange(field); err != errNotHeaderDateRange {
				// Check if error occurred
				if err != nil {
					return p.crimes, fmt.Errorf("error parsing "+
						"header date range: %s", err.Error())
				}

				// Set new skip value if not 0
				if newSkip > 0 {
					skip = newSkip
				}
			} else if footerPageNumExpr.MatchString(field) { // Check if
				// first line of footer
				skip = 5
				pageNum += 1
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
				c.DateOccurredStart = start
				c.DateOccurredEnd = end

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
				c.ReportID = reportID

				if c.DateOccurredStart == nil {
					p.logger.Printf("empty: %s\n", c)
					// TODO: Figure out why some Crimes have "empty" date_occurred ranges
				}
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
		p.parsedCrimes = true
		return p.crimes, nil

	return 0, nil

}

func (p *DrexelParser) String() string {
	return "DrexelParser"
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

// parseRangeComponent extracts a time.Time from a date range in a page header.
// The matches results slice from the headerDateRangeExpr must be provided.
//
// Along with an offset to specify if the method should extract the beginning
// or end of the header date range. offset = 0 for the beginning date in the
// range. Offset = 1 for the end of the range. Any other value will cause an
// error.
//
// If an error occurs one will be returned. Nil on success.
func (p DrexelParser) parseHeaderDate(matches []string, offset uint) (*time.Time, error) {
	// Check offset
	if offset > 1 {
		return nil, fmt.Errorf("offset argument can not be greater"+
			" than 2, was: %d", offset)
	}

	// idxOffset will be added to the indexes used to access the month, day,
	// and year in the matches slice. And it determined by the offset arg.
	idxOffset := offset * 3

	// Extract fields from matches
	monthStr := matches[idxOffset+1]
	dayStr := matches[idxOffset+2]
	yearStr := matches[idxOffset+3]

	// Parse month
	month, err := p.parseMonthAbbrv(monthStr)
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

// parseHeaderRange determines if the field provided is the report date range.
// If it is, the field is parsed and saved so the Range method can return it.
//
// The number of fields to skip after parseRange is called is returned. If 0,
// the existing skip variable should not be modified.
//
// Additionally an error is returned if one occurs, nil on success. The
// errNotHeaderDateRange error will be returned if the provided field was not
// in the header date range format.
func (p *DrexelParser) parseHeaderRange(field string) (int, error) {
	validSkip := 3
	errSkip := -1

	// Check if field is header date range
	if matches := headerDateRangeExpr.FindStringSubmatch(field); matches != nil {
		// Convert start date
		startRange, err := date.ParseDate(matches[0:3])
		if err != nil {
			return errSkip, fmt.Errorf("error converting start "+
				"header date to time.Time: %s",
				err.Error())
		}
		report.RangeStartDate = startRange

		// Convert end date
		endRange, err := date.ParseDate(matches[3:6])
		if err != nil {
			return errSkip, fmt.Errorf("error converting end "+
				"header date to time.Time: %s",
				err.Error())
		}
		report.RangeEndDate = endRange

		// Mark as parsed and skip fields
		return validSkip, nil
	}

	// If not a header date range
	return 0, errNotHeaderDateRange
}
*/
