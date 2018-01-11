package parsers

import (
	"errors"
	"strings"
	"time"

	"github.com/Noah-Huppert/crime-map/models"
)

// Parser provides methods for converting a stream of pdf text fields to Crime
// structs. This allows multiple different formats of reports to parsed.
type Parser interface {
	// Parse takes the text fields saved in the implementing struct and
	// parses them into a slice of Crime structs.
	//
	// A Report ID is provided to the Parse method. Which is the ID of the
	// Report which these crimes belong to.
	//
	// Additionally an error is returned, nil on success.
	Parse(reportID int, fields []string) ([]models.Crime, error)

	// Range returns the date range which the report covers crimes for.
	// Start time, then end time. Along with an error if one occurs, or nil
	// on success.
	//
	// The parse method is expected to be called before Range(). And an
	// error will be returned if it has not been.
	Range(fields []string) (*time.Time, *time.Time, error)

	// Count returns the number of Crime models parsed from a report. An
	// error is returned if one occurs. Nil on success.
	Count() (uint, error)
}

// determineUniversity figures out which University a crime report was
// published from. By reading in the text fields present in a report. And
// searching for the first occurrence of a university name.
//
// A models.UniversityType is returned along with an error. Which will be nil
// on success.
func determineUniversity(fields []string) (models.UniversityType, error) {
	// Attempt to find univ name in fields
	for _, field := range fields {
		// Check
		if strings.Contains(field, string(models.UniversityDrexel)) {
			// Success
			return models.UniversityDrexel, nil
		}
	}

	// If none found
	return models.UniversityErr, errors.New("error determining university," +
		" no field with university name found")
}
