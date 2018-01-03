package parsers

import (
	"github.com/Noah-Huppert/crime-map/models"
	"time"
)

// Parser provides methods for converting a stream of pdf text fields to Crime
// structs. This allows multiple different formats of reports to parsed.
type Parser interface {
	// Parse takes the provided pdf text fields and converts them into a
	// slice of Crime structs. Additionally an error is returned, nil on
	// success.
	Parse(fields []string) ([]models.Crime, error)

	// Range returns the date range which the report covers crimes for.
	// Start time, then end time. Along with an error if one occurs, or nil
	// on success.
	//
	// The parse method is expected to be called before Range(). And an
	// error will be returned if it has not been.
	Range() (*time.Time, *time.Time, error)
}
