package parsers

import (
	"fmt"

	"github.com/Noah-Huppert/crime-map/errs"
	"github.com/Noah-Huppert/crime-map/models"
)

// ParserRunner invokes a series of Parsers on an array of fields. Fields will
// parsed be parsed in order. Parsers will be invoked on these fields in a
// specified order.
//
// The next field will be parsed when any of the Parsers indicates that it
// consumed one or more fields.
//
// These Parsers can be added in the order they will be invoked.
type ParserRunner struct {
	// parsers holds the Parsers used to extract information from a list
	// of fields. Parsers in this slice will be tried one after the next
	// until a field is successfully parsed. Then the next field will be
	// parsed.
	//
	// The order of Parsers in this slice indicates precedence.
	parsers []Parser

	// crimes holds the Crime models which have been parsed so far by the
	// specified Parsers
	crimes []*models.Crime
}

// NewParserRunner creates and returns a ParserRunner reference
func NewParserRunner() *ParserRunner {
	return &ParserRunner{
		parsers: []Parser{},
		crimes:  []*models.Crime{},
	}
}

// Parse will parse the provided fields with the previously added
// Parsers. A Report model should be provided to save information about the
// report fields being parsed.
//
// The Crimes which were parsed from these fields will be returned.
// Along with any error that occurs, Nil on success.
func (r ParserRunner) Parse(report *models.Report, fields []string) ([]*models.Crime, error) {
	// crimes holds the Crime models parsed by the parsers
	crimes := []*models.Crime{}

	// crime is the Crime model currently being parsed
	crime := &models.Crime{}

	// Loop through fields and parse
	var fI uint = 0
	for fI < uint(len(fields)) {
		// Run one parser after another until a field is successfully
		// parsed. Then move to the next field until done.
		for pI, parser := range r.parsers {
			delta, err := parser.Parse(fI, fields, report, crime)

			// Check if current crime has been successfully parsed
			if err == errs.ErrCrimeParsed {
				// Add current crime to list
				crimes = append(crimes, crime)

				// Make new crime model
				crime = &models.Crime{}
			} else if err != nil {
				// If other error
				return nil, fmt.Errorf("error running %s parser "+
					"against field with index %d, err: %s", parser,
					fI, err.Error())
			}

			// Check if any fields were parsed
			if delta > 0 {
				// Increment field index
				fI += delta

				// Go to next field
				break
			} else if pI+1 == len(r.parsers) { // If last parser
				// and no delta value > 0

				// If not parsed, error
				return nil, NewErrFieldNotParsed(fI, fields[fI])
			}
		}
	}

	// Success
	return crimes, nil
}

// Add will save the provided parser for use when any fields need to
// be parsed.
func (r *ParserRunner) Add(parser Parser) {
	r.parsers = append(r.parsers, parser)
}
