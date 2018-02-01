package parsers

import (
	"fmt"
)

// ErrFieldNotParsed indicates that a field was not parsed by a parser
type ErrFieldNotParsed struct {
	// index is the slice index of the field which was not parsed
	index uint

	// field is the string that was not parsed
	field string
}

// NewErrFieldNotParsed creates a new ErrFieldNotParsed error with the specified
// index and field values
func NewErrFieldNotParsed(index uint, field string) *ErrFieldNotParsed {
	return &ErrFieldNotParsed{
		index: index,
		field: field,
	}
}

// Error implements the error interface for the ErrFieldNotParsed type
func (e *ErrFieldNotParsed) Error() string {
	return e.String()
}

// String converts the ErrFieldNotParsed struct to a string
func (e ErrFieldNotParsed) String() string {
	return fmt.Sprintf("field not parsed, index: %s, field: %s", e.index,
		e.field)
}
