package models

import (
	"database/sql"
	"fmt"

	"github.com/Noah-Huppert/crime-map/dstore"
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

// Query searches for a parse error with the same field values in the database.
// An error is returned if one occurs, or nil on success.
//
// The ParseError.ID field will be set to record the ID of the row in the
// database.
func (e *ParseError) Query() error {
	// Get database instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Query
	row := db.QueryRow("SELECT id FROM parse_errors WHERE crime_id = $1 AND "+
		"field = $2 AND original = $3 AND corrected = $4 AND err_type = $5",
		e.CrimeID, e.Field, e.Original, e.Corrected, e.ErrType)

	// Get ID
	err = row.Scan(&e.ID)

	// Check if not found
	if err == sql.ErrNoRows {
		// Return err so we can identify
		return err
	} else if err != nil {
		// Other error
		return fmt.Errorf("error querying database for ParseError model"+
			", ParseError: %s, err: %s", e, err.Error())
	}

	// Success
	return nil
}

// Insert adds a ParseError model to the database. An error is returned if one
// occurs, or nil on success.
//
// The ParseError.ID field will be set to record the ID of the newly inserted
// row.
func (e *ParseError) Insert() error {
	// Get db
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Insert
	row := db.QueryRow("INSERT INTO parse_errors (crime_id, field, original, "+
		"corrected, err_type) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		e.CrimeID, e.Field, e.Original, e.Corrected, e.ErrType)

	// Get ID
	err = row.Scan(&e.ID)
	if err != nil {
		return fmt.Errorf("error inserting ParseError model: %s",
			err.Error())
	}

	// Success
	return nil
}

// InsertIfNew will add a ParseError to the database if one with the same values
// doesn't already exist. The found / inserted row's ID will be recorded in the
// ParseError.ID field.
//
// An error will be returned if one occurs, or nil on success
func (e *ParseError) InsertIfNew() error {
	// Query
	err := e.Query()

	// Check if doesn't exist
	if err == sql.ErrNoRows {
		// Insert
		if err = e.Insert(); err != nil {
			return fmt.Errorf("error inserting non-existent "+
				"ParseError: %s", err.Error())
		}
	} else if err != nil {
		return fmt.Errorf("error querying for ParseError: %s", err.Error())
	}

	// Success
	return nil
}
