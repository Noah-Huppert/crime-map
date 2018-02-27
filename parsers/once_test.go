package parsers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Noah-Huppert/crime-map/errs"
	"github.com/Noah-Huppert/crime-map/models"
)

// IncrementParser is a parser used for testing which increments the Report and
// Crime model's ID fields. This essentailly records how many times the parser
// has been run. And is simplier than mocking a Parser.Parse method.
//
// The IncrementParser will only increment these fields if the field currently
// being parsed has the value "IncrementParser"
type IncrementParser struct{}

// incrementParserName holds the IncrementParser's name.
const incrementParserName string = "IncrementParser"

// Parse implements the Parser.Parse method for the IncrementParser
func (p IncrementParser) Parse(i uint, fields []string, report *models.Report,
	crime *models.Crime) (uint, error) {

	// Increment if field is == "IncrementParser"
	if fields[i] == p.String() {
		report.ID += 1
		crime.ID += 1

		return 1, nil
	}

	return 0, nil

}

func (p IncrementParser) String() string {
	return incrementParserName
}

// ErrorParser is a parser which always throws the errorParserErr
type ErrorParser struct{}

// errorParserName holds the ErrorParser's name
const errorParserName string = "ErrorParser"

// errorParserErr is the error which the ErrorParser always throws
var errorParserErr error = errors.New("error parser")

// Parse implements the Parser.Parse method for the ErrorParser
func (p ErrorParser) Parse(i uint, fields []string, report *models.Report,
	crime *models.Crime) (uint, error) {
	return 0, errorParserErr
}

func (p ErrorParser) String() string {
	return errorParserName
}

// TestOnceParse ensures that the OnceRunner.Parse method only runs the
// provided parser once.
func TestOnceParser(t *testing.T) {
	r := NewOnceRunner(IncrementParser{})

	// Test
	report := &models.Report{}
	crime := &models.Crime{}

	err := r.Parse(report, crime, []string{"A", "B", "IncrementParser",
		"IncrementParser", "IncrementParser"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	} else if report.ID != 1 || crime.ID != 1 {
		t.Fatalf("report or crime ID field did not equal expected "+
			"value, expected: 1, actual: report.ID: %d, crime.ID: %d",
			report.ID, crime.ID)
	}
}

// TestOnceErrors ensures that the OnceRunner.Parse method returns an error
// if one occurs with the provided parser.
func TestOnceErrors(t *testing.T) {
	r := NewOnceRunner(ErrorParser{})

	// Test
	report := &models.Report{}
	crime := &models.Crime{}

	expectedErrStr := fmt.Sprintf("error parsing field with index: 0, "+
		"field: any, err: %s", errorParserErr)
	err := r.Parse(report, crime, []string{"any", "fields"})

	if (err == nil) || (err.Error() != expectedErrStr) {
		t.Fatalf("error does not match expected: %s, actual: %s",
			expectedErrStr, err)
	}
}

// NeverParser never parses a field
type NeverParser struct{}

// neverParserName holds the NeverParser's name
const neverParserName string = "NeverParser"

// Parse implements Parser.Parse for NeverParser. It will never parse a field.
func (p NeverParser) Parse(i uint, fields []string, report *models.Report,
	crime *models.Crime) (uint, error) {
	return 0, nil
}

func (p NeverParser) String() string {
	return neverParserName
}

// TestOnceNoParseErrs ensures that OnceRunner throws an error if no fields
// are parsed
func TestOnceNoParseErrs(t *testing.T) {
	r := NewOnceRunner(NeverParser{})

	// Test
	report := &models.Report{}
	crime := &models.Crime{}

	expectedErr := errs.ErrNotParsed
	err := r.Parse(report, crime, []string{"none", "will", "parse"})

	if (err == nil) || (err != expectedErr) {
		t.Fatalf("error did not match expected: %s, actual: %s",
			expectedErr, err)
	}
}
