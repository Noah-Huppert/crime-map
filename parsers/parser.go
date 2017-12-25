package parsers

import (
	"github.com/Noah-Huppert/crime-map/models"
)

// Parser provides methods for converting a stream of pdf text fields to Crime
// structs. This allows multiple different formats of reports to parsed.
type Parser interface {
	// Parse takes the provided pdf text fields and converts them into a
	// slice of Crime structs. Additionally an error is returned, nil on
	// success.
	Parse(fields []string) ([]models.Crime, error)
}
