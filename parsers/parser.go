package parsers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Noah-Huppert/crime-map/models"
)

// Parser takes a single text field and extracts information related to
// the general crime report document or a single reported crime.
//
// Parser implementations should each only handle 1 type of data field. Multiple
// Parsers can then be run together to handle every type of text field a crime
// or report may contain.
type Parser interface {
	// Parsers must implement String() which will return a identifying
	// name of the parser to use in errors and various state information
	fmt.Stringer

	// Parse extracts information from a text field contained in a crime
	// report pdf.
	//
	// The index of the field to parse is provided by the `i` argument. The
	// slice of fields to access is provided by the `fields` argument.
	//
	// Information which is parsed should be saved in the models.Report and
	// models.Crime pointers.
	//
	// The number of fields which were parsed should be returned. 0 if no
	// fields were parsed.
	//
	// Finally an error is returned if one occurs. Nil on success.
	//
	// To indicate the provided crime is fully parsed return the
	// ErrCrimeParsed value. This will indicate to the caller that they
	// should add the currently provided crime to a list of parsed crimes.
	// And provide a new crime model reference for the next invocation.
	Parse(i uint, fields []string, report *models.Report, crime *models.Crime) (uint, error)
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
