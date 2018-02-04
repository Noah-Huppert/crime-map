package parsers

import (
	"testing"
	"time"

	"github.com/Noah-Huppert/crime-map/models"
)

// TestDateRangeParserNotRange ensures that the DateRangeParser does not parse
// a field which is not a date range.
func TestDateRangeParserNotRange(t *testing.T) {
	p := DateRangeParser{}

	// Test number of parsed fields is 0
	report := &models.Report{}
	numParsed, err := p.Parse(0, []string{"not a date range"}, report, nil)

	if err != nil {
		t.Fatalf("error running DateRangeParser: %s", err.Error())
	} else if numParsed != 0 {
		t.Fatalf("number of parsed fields not 0, was: %d", numParsed)
	}
}

// expectedStartErr is the error that is expected to be thrown in
// TestDateRangeParserErrsStart.
const expectedStartErr string = "error converting start header date to " +
	"time.Time: error parsing month abbreviation into " +
	"month number: error parsing month abbreviation, " +
	"unknown value: Man"

// TestDateRangeParserErrsStart ensures that the DateRangeParser throws an error
// when there is a formatting issue with the start of the date range.
func TestDateRangeParserErrsStart(t *testing.T) {
	p := DateRangeParser{}

	// Test error gets thrown if the date range start is mis-formatted
	report := &models.Report{}
	_, err := p.Parse(0, []string{"From Man 13, 2016 to Jan 13, 2017."}, report, nil)

	if err != nil {
		// Check error matches expected
		if expectedStartErr != err.Error() {
			t.Fatalf("did not receive expected error, \n"+
				"expected: \"%s\", \n"+
				"actual  : \"%s\"",
				expectedStartErr, err.Error())
		}
	} else {
		t.Fatalf("no error")
	}
}

// expectedEndErr is the error that is expected to be thrown in
// TestDateRangeParserErrsEnd.
const expectedEndErr string = "error converting end header date to " +
	"time.Time: error parsing month abbreviation into " +
	"month number: error parsing month abbreviation, " +
	"unknown value: San"

// TestDateRangeParserErrsEnd ensures that the DateRangeParser throws an error
// when there is a formatting issue with the end of the date range.
func TestDateRangeParserErrsEnd(t *testing.T) {
	p := DateRangeParser{}

	// Test error gets thrown if the date range start is mis-formatted
	report := &models.Report{}
	_, err := p.Parse(0, []string{"From Jan 13, 2016 to San 13, 2017."}, report, nil)

	if err != nil {
		// Check error matches expected
		if expectedEndErr != err.Error() {
			t.Fatalf("did not receive expected error, \n"+
				"expected: \"%s\", \n"+
				"actual  : \"%s\"",
				expectedEndErr, err.Error())
		}
	} else {
		t.Fatalf("no error")
	}
}

// expectedNumParsed is the expected number of fields the DateRangeParser should
// indicate it parsed.
const expectedNumParsed int = 3

// expectedRangeStart is the expected range start date that the DateRangeParser
// should parse and save in the provided Report model.
var expectedRangeStart = time.Date(2016, time.Month(1), 13, 0, 0, 0,
	0, time.UTC)

// expectedRangeEnd is the expected range end date that the DateRangeParser
// should parse and save in the provided Report model.
var expectedRangeEnd = time.Date(2017, time.Month(1), 13, 0, 0, 0, 0,
	time.UTC)

// TestDateRangeParserSavesInfo ensures that the DateRangeParser saves the
// correct information in the provided Report. And that it specifies the correct
// number of fields have been parsed.
func TestDateRangeParserSavesInfo(t *testing.T) {
	p := DateRangeParser{}

	// Test saves info in report model
	report := &models.Report{}
	numParsed, err := p.Parse(0, []string{"From Jan 13, 2016 to Jan 13, 2017."}, report, nil)

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	} else if numParsed != 3 {
		t.Fatalf("incorrect number of fields parsed, expected: 3, "+
			"actual: %d", numParsed)
	} else if *report.RangeStartDate != expectedRangeStart {
		t.Fatalf("report.RangeStartDate does not matched expected value,"+
			"expected: %s, actual: %s", expectedRangeStart,
			report.RangeStartDate)
	} else if *report.RangeEndDate != expectedRangeEnd {
		t.Fatalf("report.RangeEndDate does not matched expected value,"+
			"expected: %s, actual: %s", expectedRangeEnd,
			report.RangeEndDate)
	}
}

// TestDateRangeParserString ensures the DateRangeParser.String method returns
// the correct value.
func TestDateRangeParserString(t *testing.T) {
	p := DateRangeParser{}

	// Test DateRangeParser.String
	actual := p.String()

	if actual != DateRangeParserName {
		t.Fatalf("String() value does not match expected, expected: %s"+
			", actual: %s", DateRangeParserName, actual)
	}
}
