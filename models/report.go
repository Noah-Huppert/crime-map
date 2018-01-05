package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Noah-Huppert/crime-map/dstore"
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

// Report holds information about documents parsed by crime-map to extract
// Crime models. These documents are typically published as an obligation to
// the Clery Act. And hold multiple individual crime reports.
type Report struct {
	// ID is the unique report identifier
	ID int

	// ParsedOn indicates the date and time the report was processed
	ParsedOn *time.Time

	// ParseSuccess indicates if the report was successfully parsed
	ParseSuccess bool

	// University indicates which institution published the crime report
	// document
	University UniversityType

	// RangeStartDate indicates the start of the date range crimes were reported
	// for
	RangeStartDate *time.Time

	// RangeEndDate indicates the end of the date range crimes were reported
	// for
	RangeEndDate *time.Time

	// Pages holds the number of pages the document had
	Pages uint
}

// NewReport will create a new Report model.
func NewReport(univ UniversityType, parsedOn *time.Time, start *time.Time,
	end *time.Time, pages uint) *Report {
	return &Report{
		ParsedOn:       parsedOn,
		ParseSuccess:   false,
		University:     univ,
		RangeStartDate: start,
		RangeEndDate:   end,
		Pages:          pages,
	}
}

// NewReportFromRow creates a new Report model from a database row. This row
// should be from a query which selects the id, parsed_on, parse_success,
// university, range, and pages fields. Additionally an error is returned if
// one occurs, nil on success.
func NewReportFromRow(rows *sql.Rows) (*Report, error) {
	// Scan
	r := &Report{}
	var dRange interface{}
	err := rows.Scan(&r.ID, &r.ParsedOn, &r.ParseSuccess, &r.University,
		&dRange, &r.Pages)

	if err != nil {
		return nil, fmt.Errorf("error parsing Report from database row"+
			": %s", err.Error())
	}

	// Success
	return r, nil
}

// String encodes the Report into string form
func (r Report) String() string {
	return fmt.Sprintf("ID: %d\n"+
		"ParsedOn: %s\n"+
		"ParseSuccess: %t\n"+
		"University: %s\n"+
		"Range: [%s, %s]\n"+
		"Pages: %d",
		r.ID, r.ParsedOn, r.ParseSuccess, r.University,
		r.RangeStartDate, r.RangeEndDate, r.Pages)
}

// Query attempts to find a Report with the same parse_success, university,
// range, and pages field values. The parsed_on field is left out of the query.
//
// It populates the Report.ID field with the database row's ID. An error is
// returned if one occurs, or nil on success.
func (r *Report) Query() error {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Query
	row := db.QueryRow("SELECT id FROM reports WHERE parse_success=$1 AND "+
		"university=$2 AND range=tstzrange($3, $4, '()') AND pages=$5",
		r.ParseSuccess, r.University, r.RangeStartDate, r.RangeEndDate,
		r.Pages)

	// Get ID
	err = row.Scan(&r.ID)

	// Check if no rows
	if err == sql.ErrNoRows {
		// Return error so we can identify
		return err
	} else if err != nil {
		// Other error
		return fmt.Errorf("error querying for Report: %s", err.Error())
	}

	// Success
	return nil
}

// Insert adds a Report model to the database. An error is returned if one
// occurs, or nil on success.
//
// The ID of the newly inserted row will be saved in the Report.ID field.
func (r *Report) Insert() error {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Insert
	row := db.QueryRow("INSERT INTO reports (parsed_on, parse_success, "+
		"university, range, pages) VALUES ($1, $2, $3, "+
		"tstzrange($4, $5, '()'), $6) RETURNING id",
		r.ParsedOn, r.ParseSuccess, r.University, r.RangeStartDate,
		r.RangeEndDate, r.Pages)

	// Get ID
	err = row.Scan(&r.ID)
	if err != nil {
		return fmt.Errorf("error inserting Report model: %s",
			err.Error())
	}

	// Success
	return nil
}

// InsertIfNew adds a Report model to the database if one with existing values
// does not exist yet. The ID of the queried/inserted row is saved in the
// Report.ID field. An error is returned if one occurs, nil on success.
func (r *Report) InsertIfNew() error {
	// Query
	err := r.Query()

	// If doesn't exist
	if err == sql.ErrNoRows {
		// Insert
		if err = r.Insert(); err != nil {
			return fmt.Errorf("error inserting non-existing "+
				"Report model: %s", err.Error())
		}
	} else if err != nil {
		// Other error
		return fmt.Errorf("error querying for existence of Report"+
			" model: %s", err.Error())
	}

	// Success
	return nil
}

// QueryAllReports finds all Report models from the database. And returns them
// with their Report.ID fields populated. Additionally an error is returned if
// one occurs. Nil on success.
func QueryAllReports() ([]*Report, error) {
	reports := []*Report{}

	// Get db
	db, err := dstore.NewDB()
	if err != nil {
		return reports, fmt.Errorf("error retrieving database "+
			"instance: %s", err.Error())
	}

	// Query
	rows, err := db.Query("SELECT id, parsed_on, parse_success, university" +
		", range, pages FROM reports ORDER BY parsed_on DESC")

	// Parse
	for rows.Next() {
		report, err := NewReportFromRow(rows)
		if err != nil {
			return reports, fmt.Errorf("error parsing report row: %s",
				err.Error())
		}

		reports = append(reports, report)
	}

	// Close query
	if err = rows.Close(); err != nil {
		return reports, fmt.Errorf("error closing reports query: %s",
			err.Error())
	}

	// Success
	return reports, nil
}
