package parsers

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/Noah-Huppert/crime-map/errs"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/test"
)

// SampleIDA is the value which SampleParserA will set the crime.ID field to
const SampleIDA int = 5

// SampleIDB is the value which SampleParserB will set the crime.ID field to
const SampleIDB int = 11

// SampleParserA is a Parser which parse fields with the value "A". It will set
// the Crime.ID field to SampleIDA. And then returns errs.ErrCrimeParsed.
type SampleParserA struct{}

const SampleParserAName string = "SampleParserA"

func (p SampleParserA) String() string {
	return SampleParserAName
}

func (p SampleParserA) Parse(i uint, fields []string, report *models.Report,
	crime *models.Crime) (uint, error) {
	// Check if field is "A"
	if fields[i] == "A" {
		// Set crime.ID to SampleIDA
		crime.ID = SampleIDA

		// Indicate crime parsed
		return 1, errs.ErrCrimeParsed
	}

	// Otherwise don't parse
	return 0, nil
}

// SampleParserB is a Parser which parse fields with the value "B". It will set
// the Crime.ID field to SampleIDB. And then returns errs.ErrCrimeParsed.
type SampleParserB struct{}

const SampleParserBName string = "SampleParserB"

func (p SampleParserB) String() string {
	return SampleParserBName
}

func (p SampleParserB) Parse(i uint, fields []string, report *models.Report,
	crime *models.Crime) (uint, error) {
	// Check if field is "B"
	if fields[i] == "B" {
		// Set crime.ID to SampleIDA
		crime.ID = SampleIDB

		// Indicate crime parsed
		return 1, errs.ErrCrimeParsed
	}

	// Otherwise don't parse
	return 0, nil
}

// SampleErrParser implements the Parser interface by always throwing an error
// when a field is given to parse
type SampleErrParser struct{}

// SampleErrParserName is the name of the SampleErrParser
const SampleErrParserName string = "SampleErrParser"

// SampleParserErr is the error returned by the SampleErrParser
var SampleParserErr = errors.New("sample error")

func (p SampleErrParser) String() string {
	return SampleErrParserName
}

func (p SampleErrParser) Parse(i uint, fields []string, report *models.Report,
	crime *models.Crime) (uint, error) {

	return 0, SampleParserErr
}

// TestParserRunnerAddParse ensures the ParserRunner.Add & .Parse methods
// work as expected
func TestParserRunnerAddParse(t *testing.T) {
	// Make runner
	runner := NewParserRunner()

	// Add mock Parsers
	runner.Add(SampleParserA{})
	runner.Add(SampleParserB{})

	// Parse
	report := &models.Report{}
	crimes, err := runner.Parse(report, []string{"B", "A", "B", "B", "A"})
	if err != nil {
		t.Fatalf("error parsing fields with ParserRunner: %s",
			err.Error())
	}

	// Check
	test.CrimesSlicesEq(t, []*models.Crime{
		&models.Crime{ID: SampleIDB},
		&models.Crime{ID: SampleIDA},
		&models.Crime{ID: SampleIDB},
		&models.Crime{ID: SampleIDB},
		&models.Crime{ID: SampleIDA},
	}, crimes)
}

// TestParserRunnerNoParserErr ensures that ParserRunner.Parse throws an error
// when a field is encountered that does not have an associated parser.
func TestParserRunnerNoParserErr(t *testing.T) {
	// Make runner
	runner := NewParserRunner()

	// Add Mock Parsers
	runner.Add(SampleParserA{})
	runner.Add(SampleParserB{})

	// Parse
	report := &models.Report{}
	_, err := runner.Parse(report, []string{"C", "A"})

	// Save error to check later
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	// Ensure match
	matched, err := regexp.MatchString("no parser processed "+
		"field with index [0-9]*, field: .*", errStr)
	if err != nil {
		t.Fatalf("error matching ParserRunner.Parse error: %s",
			err.Error())
	} else if !matched {
		t.Fatalf("ParserRunner.Parse error does not match expected "+
			"pattern, actual: \"%s\"", errStr)
	}
}

// TestParserRunnerErr tests the ParserRunner.Parse method when a parser throws
// an error
func TestParserRunnerErr(t *testing.T) {
	// Make runner
	runner := NewParserRunner()

	// Add mock parsers
	runner.Add(SampleParserA{})
	runner.Add(SampleErrParser{})

	// Parse
	report := &models.Report{}
	_, err := runner.Parse(report, []string{"C", "A"})

	// Save error to check later
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	// Ensure match
	matched, err := regexp.MatchString(fmt.Sprintf("error running %s "+
		"parser against field with index 0, err: %s",
		SampleErrParserName, SampleParserErr), errStr)
	if err != nil {
		t.Fatalf("error matching ParserRunner.Parser error: %s",
			err.Error())
	} else if !matched {
		t.Fatalf("ParserRunner.Parse error does not match expected "+
			"pattern, actual: \"%s\"", errStr)
	}
}
