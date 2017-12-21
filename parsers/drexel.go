package parsers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Noah-Huppert/crime-map/crime"
)

// startOfCrimeToken holds the pdf text field value which appears at the start
// of a crime's fields in a report
const startOfCrimeToken string = "Location :"

// firstHeaderSubToken holds the pdf text field value which appears first in a
// page header
const firstHeaderSubToken string = "From"

// lastHeaderToken holds the pdf text field value which appears last in a page
// header
const lastHeaderToken string = "Student Right To Know Case Log Daily Report"

// ignoredFields holds all the pdf text field values that can safely be skipped
// while parsing a pdf crime log
var ignoredFields []string = []string{
	startOfCrimeToken,
	" Report #:",
	"Date and Time Occurred From - Occurred To:",
	"Print Date and Time",
	"Incident\\s\\",
	"Date Reported:",
	"Disposition:",
	"Synopsis:",
	"at",
}

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
	// Loop through crime fields and transform into Crime structs
	crimeFields := p.GroupFieldsByCrime(fields)
	crimes := []crime.Crime{}

	for _, fields := range crimeFields {
		// Check at least 5 fields provided
		numFields := len(fields)
		if numFields < 5 {
			return crimes, fmt.Errorf("error parsing crime: not enough "+
				"fields, found %d, needs 5, %s", len(fields),
				fields)
		}

		// Set fields
		crime := crime.Crime{}

		// TODO: Parse into more meaningful data types
		crime.DateReported = fields[0]
		crime.Location = fields[1]
		crime.ReportID = fields[2]
		crime.Incidents = []string{fields[3]}
		crime.DateOccurred = fields[4]

		// If description provided
		if numFields > 5 {
			crime.Descriptions = fields[4 : numFields-2]
		}

		crime.Remediation = fields[numFields-1]

		crimes = append(crimes, crime)
	}

	return crimes, nil

}

// IsIgnoredField returns a boolean indicating whether or not the provided
// field should be ignored
func (p DrexelParser) IsIgnoredField(field string) bool {
	// Loop through all ignored fields
	for _, val := range ignoredFields {
		// If match
		if field == val {
			return true
		}
	}

	// No match
	return false
}

// GroupFieldsByCrime separates fields into common arrays based on the crime
// they belong to. And returns a list, of list of fields.
func (p DrexelParser) GroupFieldsByCrime(fields []string) [][]string {
	// grouped holds fields, grouped by which crime they belong to, and
	// will be returned at the end of the fn
	grouped := [][]string{}

	// inHeader indicates if the program is currently parsing fields
	// which are in the pdf reports header
	inHeader := false

	// crimeFields holds the fields of the crime which is currently being
	// grouped. The contents of which will be added to the grouped var
	// every time a new crime starts
	crimeFields := []string{}

	crimesI := 0

	// Loop through fields
	for _, field := range fields {
		// Start of a new crime
		if field == startOfCrimeToken {
			// Add to grouped var if not empty
			if len(crimeFields) > 0 {
				grouped = append(grouped, crimeFields)
			}

			// Reset crimeFields
			crimeFields = []string{}

			crimesI++
		} else {
			// If midway through parsing a crime's fields OR in
			// between crimes

			// Check if page number
			// By trying to parse the current field into an int
			if _, err := strconv.Atoi(field); err == nil {
				inHeader = true
				continue
			}

			// Check if field is last item in a page header
			if field == lastHeaderToken {
				// Mark out of header
				inHeader = false
				continue
			}

			// Check if field is date range report covers, first
			// item in header
			// In form "From xxx to xxx", so just check first word
			split := strings.Split(field, " ")
			if len(split) > 0 && split[0] == firstHeaderSubToken {
				inHeader = true
				continue
			} else if inHeader {
				continue
			}

			// Check if ignored
			if p.IsIgnoredField(field) {
				continue
			}

			// If total count of crimes at bottom of report
			if field == fmt.Sprintf(" %d", crimesI+1) {
				continue
			}

			// If normal text field, Add to crime fields var
			crimeFields = append(crimeFields, field)
		}
	}

	// Add last crime fields result if not empty
	if len(crimeFields) > 0 {
		grouped = append(grouped, crimeFields)
	}

	// All done
	return grouped
}
