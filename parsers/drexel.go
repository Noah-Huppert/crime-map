package parsers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Noah-Huppert/crime-map/crime"
)

// Header constants
// headerDateRangeExpr is the regexp used to match a report header's first line
var headerDateRangeExpr *regexp.Regexp = regexp.MustCompile("^From [A-Z][a-z]+ [0-9]{1,2}, [0-9]{4} to [A-Z][a-z]+ [0-9]{1,2}, [0-9]{4}\\.$")

// headerUnivName is the name of university which appears in the report which
// appears in the header
const headerUnivName string = "Drexel University"

// headerReportName is the name of the report which appears in the header
const headerReportName string = "Public Safety"

// lastHeaderToken holds the pdf text field value which appears last in a page
// header
const headerTitle string = "Student Right To Know Case Log Daily Report"

// Footer constants
// footerPageNumExpr holds the regexp used to match page numbers at the bottom
// of the header
var footerPageNumExpr *regexp.Regexp = regexp.MustCompile("^[0-9]+$")

// footerPageNumLabel holds the text which appears at the bottom of a footer
const footerPageNumLabel string = "Page No."

// footerTimeExpr is the regexp used to match the time in the footer of a
// report
var footerTimeExpr *regexp.Regexp = regexp.MustCompile("^[0-9]{2}:[0-9]{2}:[0-9]{2}$")

// footerDateExpr is the regexp used to match the date in the footer of a report
var footerDateExpr *regexp.Regexp = regexp.MustCompile("^([0-9]{2}/){2}[0-9]{4}$")

// footerPrintLabel is the date and time label which appears in the footer of a
// report
const footerPrintLabel string = "Print Date and Time"

// footerPrintLabel2 is the second part of the time and date label which
// appears in the footer of a report
const footerPrintLabel2 string = "at"

// Field labels
const fieldLabelReported string = "Date Reported:"
const fieldLabelID string = " Report #:"
const fieldLabelLoc string = "Location :"
const fieldLabelIncidents string = "Incident(s):"
const fieldLabelOccurred string = "Date and Time Occurred From - Occurred To:"
const fieldLabelDesc string = "Synopsis:"
const fieldLabelFix string = "Disposition:"

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
	crimes := []crime.Crime{}

	// globber holds fields that are currently being parsed, for context
	globber := []string{}
	// footerGlobber holds fields that might potentially be part of a page
	// footer. Fields only get added if len(footerGlobber) > 0 or a field
	// matches footerPageNumExpr
	footerGlobber := []string{}

	globLength := 0
	footerGlobLength := 0

	// Loop through fields
	for _, field := range fields {
		// TODO fix footer being in middle of other fields issue

		// Check if field is a number
		if footerGlobLength == 0 &&
			footerPageNumExpr.MatchString(field) {
			// Reset and add to footer globber
			footerGlobber = []string{field}
		} else if (footerGlobLength == 1 &&
			field == footerPageNumLabel) ||
			(footerGlobLength == 2 &&
				footerTimeExpr.MatchString(field)) ||
			(footerGlobLength == 3 &&
				footerDateExpr.MatchString(field)) ||
			(footerGlobLength == 4 &&
				field == footerPrintLabel) ||
			(footerGlobLength == 5 &&
				field == footerPrintLabel2) { // Check if footer globber is filling
			footerGlobber = append(footerGlobber, field)
		} else {
			// Add to globber
			globber = append(globber, field)
		}

		footerGlobLength = len(footerGlobber)
		globLength = len(globber)

		// Check if globber contains a header
		fmt.Printf("%d || %d | ", globLength, footerGlobLength)
		fmt.Println(field)
		if globLength == 4 &&
			headerDateRangeExpr.MatchString(globber[0]) &&
			globber[1] == headerUnivName &&
			globber[2] == headerReportName &&
			globber[3] == headerTitle {
			// If header, just reset globber
			fmt.Println("===== ^HEADER^ =====\n")
			globber = []string{}
			footerGlobber = []string{}
		} else if footerGlobLength == 6 && // Check if globber contains footer
			footerGlobber[1] == footerPageNumLabel &&
			footerTimeExpr.MatchString(footerGlobber[2]) &&
			footerDateExpr.MatchString(footerGlobber[3]) &&
			footerGlobber[4] == footerPrintLabel &&
			footerGlobber[5] == footerPrintLabel2 {
			// If footer, just reset globber
			fmt.Println("===== ^FOOTER^ =====\n")
			footerGlobber = []string{}

			if globLength >= 6 {
				globber = globber[:globLength-6]
			} else {
				globber = []string{}
			}
		} else if globLength == 7 && // Check if globber contains first part of crime data
			globber[0] == fieldLabelReported &&
			globber[1] == fieldLabelID &&
			globber[2] == fieldLabelLoc {
			// First first glob of a crime's info
			// reported, location, report #, incidents
			dateReported := globber[3]
			location := globber[4]
			reportID := globber[5]
			incidents := globber[6]

			fmt.Println("===== ^GLOB1^ =====")
			fmt.Printf("    Date reported: %s\n"+
				"    Location: %s\n"+
				"    Report ID: %s\n"+
				"    Incidents: %s\n\n",
				dateReported,
				location,
				reportID,
				incidents)
			globber = []string{}
		} else if globLength >= 4 && // Check if globber contains second part of crime data
			globber[0] == fieldLabelIncidents &&
			globber[2] == fieldLabelOccurred &&
			globber[3] == fieldLabelDesc &&
			globber[globLength-2] == fieldLabelFix {
			// Second glob of crime's info

			dateOccurred := globber[1]
			synopsis := []string{}

			// See if any synopsis provided
			if globLength > 6 {
				synopsis = globber[4:6]
			}

			disposition := globber[globLength-1]

			fmt.Println("===== ^GLOB2^ =====")
			fmt.Printf("    Date Occurred: %s\n"+
				"    Synopsis: %s\n"+
				"    Disposition: %s\n\n",
				dateOccurred,
				strings.Join(synopsis, ","),
				disposition)
			globber = []string{}
		}
	}

	return crimes, nil

}
