package parsers

import (
	"testing"

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
