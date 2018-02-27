package parsers

import (
	"fmt"

	"github.com/Noah-Huppert/crime-map/errs"
	"github.com/Noah-Huppert/crime-map/models"
)

// OnceRunner runs the provided Parser until it parses one or more fields
type OnceRunner struct {
	// runner will be used to run the provided Parser on the provided fields
	parser Parser
}

// NewOnceRunner will create a new OnceRunner for the provided Parser
func NewOnceRunner(parser Parser) *OnceRunner {
	return &OnceRunner{parser: parser}
}

// Parse will run the specified Parser on the provided fields. A slice of
// Crime models will be returned.
//
// Report and Crime models should be provided to save information about the
// parsing operation.

// Additionally an error will be returned if one occurs, nil on success.
func (r OnceRunner) Parse(report *models.Report, crime *models.Crime, fields []string) error {
	// Loop through fields
	for i, field := range fields {
		// Attempt to parse field
		count, err := r.parser.Parse(uint(i), fields, report, crime)

		// On error
		if (err != nil) && (err != errs.ErrCrimeParsed) {
			return fmt.Errorf("error parsing field with index"+
				": %d, field: %s, err: %s", i, field, err.Error())
		}

		// Determine if any fields parsed
		if count > 0 {
			// Success
			return nil
		}
	}

	// If couldn't parse
	return errs.ErrNotParsed
}
