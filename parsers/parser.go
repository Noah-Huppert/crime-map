package parsers

import (
	"errors"
	"strings"
	"time"

	"github.com/Noah-Huppert/crime-map/models"
)

// Parser provides methods for converting a stream of pdf text fields to Crime
// structs. This allows multiple different formats of reports to parsed.
/*
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
*/
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

// ErrCrimeParsed indicates that the provided crime has been completely
// parsed during the invocation of the parser.
//
// If this error is received: Append the crime you are currently providing
// to a list of parsed crimes. Then create a new empty crime model. And provide
// this on the next invocation of the parser.
var ErrCrimeParsed string = errors.New("crime has been successfully parsed")

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
