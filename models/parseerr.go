package models

import (
	"fmt"
)

const (
	// TypeBadRangeEnd signifies that a date range's end date occurred
	// before a range's start date.
	TypeBadRangeEnd string = "BAD_RANGE_END"
)

// ParseError structs holds details about errors which occur while parsing
// crimes from reports. And the actions that take place to fix them.
//
// This information is recorded just in case the crime was fixed incorrectly.
type ParseError struct {
	// ID is the unique identifier of the parse error
	ID int

	// CrimeID holds the ID of the crime which was corrected
	CrimeID int

	// Field holds the name of the crime field which was corrected
	Field string

	// Original holds the value of the field before it was corrected
	Original string

	// Corrected holds the value of the field after it was corrected
	Corrected string

	// ErrType holds a computer identifiable value for the error that
	// occurred
	ErrType string
}

// String converts a ParseError into a string
func (e ParseError) String() string {
	return fmt.Sprintf("ID: %d\n"+
		"Crime ID: %d\n"+
		"Field: %s\n"+
		"Original: %s\n"+
		"Corrected: %s\n"+
		"ErrType: %s",
		e.ID, e.CrimeID, e.Field, e.Original, e.Corrected, e.ErrType)
}

// StringParseErrors converts a slice of Parse Errors to a slice of strings
func StringParseErrors(errs []ParseError) []string {
	s := []string{}

	// Loop through errs
	for i, e := range errs {
		str := e.String()

		// Check not last item
		if i+1 != len(errs) {
			str += "\n"
		}

		// Convert to string
		s = append(s, str)
	}

	return s
}
