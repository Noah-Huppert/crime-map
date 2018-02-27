package models

import (
	"fmt"
)

// UniversityType is a string type alias, used to represent the valid
// university's a report can be published from.
type UniversityType string

const (
	// UniversityDrexel indicates that a report was published by Drexel
	UniversityDrexel UniversityType = "Drexel University"

	// UniversityErr indicates that a report was provided with an invalid
	// value
	UniversityErr UniversityType = "Err"
)

// NewUniversityType constructs a new valid UniversityType from a raw string.
// An error is returned if one occurs, or nil on success.
func NewUniversityType(raw string) (UniversityType, error) {
	if raw == string(UniversityDrexel) {
		return UniversityDrexel, nil
	} else {
		return UniversityErr, fmt.Errorf("error creating UniversityType"+
			" from value, invalid: %s", raw)
	}
}
