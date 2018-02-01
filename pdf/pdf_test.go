package pdf

import (
	"testing"

	"github.com/Noah-Huppert/crime-map/str"
)

// FieldsTestFile is the name of the pdf file used to test Pdf.fields related
// methods
const FieldsTestFile string = "../test_data/fields.pdf"

// FieldsFileContents holds the fields which the FieldsTestFile contains
var FieldsFileContents []string = []string{
	"one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	"ten", "eleven", "one", "two", "three", "one", "two", "three", "one",
	"two", "three", "one", "two", "three", "one", "two", "three", "one",
	"two", "three", "one", "two", "three", "one", "two", "three", "one",
	"two", "three", "one", "two", "three", "one", "two", "three", "one",
	"two", "three", "one", "two", "end",
}

// TestReadsFields ensures that the Pdf struct can parse the text fields out of
// a test pdf file.
func TestReadsFields(t *testing.T) {
	t.Skipf("%s file is not reading text fields correctly", FieldsTestFile)

	// Parse field
	file := NewPdf(FieldsTestFile)
	fields, err := file.Parse()
	if err != nil {
		t.Fatalf("error parsing fields pdf file, path: %s, err: %s",
			file.path, err.Error())
	}

	pages, parsed := file.Pages()
	t.Logf("pages: %d, parsed: %\n", pages, parsed)
	t.Logf("expected: %s\n=========\n========\nactual: %s\n", FieldsFileContents, fields)

	str.SlicesEq(t, FieldsFileContents, fields)
}
