package parsers

import (
	"testing"
	"time"

	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/test"
)

// fieldsTxtFile holds the path to the text file used to test Parsers. This
// file holds one field value per line. And is parsed into a []string, to act
// like a pdf.Pdf.Parse() output.
const fieldsTxtFile string = "test_data/fields.out"

// fakeReportID is the Report.ID crimes parsed during testing will be assigned
// to
const fakeReportID int = 5

// expectedCrimes holds the list of Crime objects expected to be parsed out of
// the fieldsTxtFile fields
var expectedCrimes []*models.Crime = []*models.Crime{}

// reportRangeStart is the start date the test report fields indicate they cover
// in the header
var reportRangeStart time.Time = time.Date(2017, time.Month(10), 14, 0, 0, 0, 0, time.UTC)

// reportRangeEnd is the end date the test report fields indicate they cover
// in the header
var reportRangeStart time.Time = time.Date(2017, time.Month(12), 14, 0, 0, 0, 0, time.UTC)

// AssertParser will ensure that the provided parser works correctly. By
// running on a test set of fields and comparing the result to the expected
// var
func AssertParser(t *testing.T, p Parser) {
	// Read fields from file
	fields, err := test.ReadFile(fieldsTxtFile)
	if err != nil {
		t.Fatalf("error reading test fields file: %s", err.Error())
	}

	// Parse fields into crimes
	actual, err := p.Parse(fields, fakeReportID)
	if err != nil {
		t.Fatalf("error parsing for fields: %s", err.Error())
	}

	// Compare crimes
	test.CrimesSlicesEq(t, expectedCrimes, actual)

	// Compare range
	start, end, err := p.Range(fields)
	if err != nil {
		t.Fatalf("error parsing report range: %s", err.Error())
	}

	if *start != reportRangeStart {
		t.Fatalf("expected report start range: %s, got: %s",
			reportRangeStart, start)
	}

	if *end != reportRangeEnd {
		t.Fatalf("expected report end range: %s, got: %s",
			reportRangeEnd, end)
	}

	// TODO: Finish full integration test
}
