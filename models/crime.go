package models

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"strings"
	"time"

	"github.com/Noah-Huppert/crime-map/date"
	"github.com/Noah-Huppert/crime-map/dstore"
)

// OrderByType is the type alias used to represent the field used to order rows
// by in the QueryAll method
type OrderByType string

const (
	// OrderByReported indicates that the QueryAll function should order
	// results by the date_reported field
	OrderByReported OrderByType = "date_reported"

	// OrderByOccurred indicates that the QueryAll function should order
	// results by the date_occurred field
	OrderByOccurred OrderByType = "date_occurred"

	// OrderByErr indicates that the created OrderByType had an invalid
	// string value
	OrderByErr OrderByType = "invalid"
)

// NewOrderByType creates an OrderByType with the specified string value. An
// error is returned if the provided string value is not a valid OrderByType.
// Or nil on success.
func NewOrderByType(val string) (OrderByType, error) {
	if val == string(OrderByReported) {
		return OrderByReported, nil
	} else if val == string(OrderByOccurred) {
		return OrderByOccurred, nil
	} else {
		return OrderByErr, fmt.Errorf("invalid OrderByType value: %s",
			val)
	}
}

// Crime structs hold information about criminal activity reported by Clery
// act reports
type Crime struct {
	// ID is a unique identifier
	ID int

	// ReportID is the unique identifier of the Report the crime was parsed
	// from
	ReportID int

	// Page indicates which page of the report a crime was reported on
	Page int

	// DateReported records when the criminal activity was disclosed to the
	// police
	DateReported time.Time

	// DateOccurredStart records when the criminal activity started taking
	// place
	DateOccurredStart *time.Time

	// DateOccurredEnd records when the criminal activity stopped taking
	// place
	DateOccurredEnd *time.Time

	// ReportSuperID is the first portion of police report ID associated
	// with the reported crime.
	ReportSuperID uint

	// ReportSubID is the second portion of the police report ID associated
	// with the reported crime
	ReportSubID uint

	// GeoLocID is the unique ID of the Geo entry which holds the geographically
	// encoded location in lat long form
	GeoLocID int

	// Incidents holds the official classifications of the criminal
	// activity which took place
	Incidents pq.StringArray `gorm:"type:text[]"`

	// Descriptions holds any details about the specific incidents which
	// took place
	Descriptions pq.StringArray `gorm:"type:text[]"`

	// Remediation is the action taken by the institution who reported the
	// crime to deal with the criminal activity
	Remediation string

	// ParseErrors holds any errors that occur while parsing the crime.
	// These will be saved in other db tables depending on their types.
	//
	// This field is used internally only. Not serialized and sent as
	// part of any API responses.
	ParseErrors []ParseError `json:"-"`
}

// NewCrime creates a new Crime model from a database query sql.Rows
// result set. This query should select the id, report_id, page, university,
// date_reported, date_occurred, report_super_id, report_sub_id, geo_loc_id,
// incidents, descriptions, and remediations fields.
//
// An Crime instance and error is returned. Nil on success.
func NewCrime(rows *sql.Rows) (*Crime, error) {
	crime := &Crime{}

	// Parse
	var rangeStr string

	if err := rows.Scan(&crime.ID, &crime.ReportID, &crime.Page,
		&crime.DateReported, &rangeStr, &crime.ReportSuperID,
		&crime.ReportSubID, &crime.GeoLocID, &crime.Incidents,
		&crime.Descriptions, &crime.Remediation); err != nil {
		return nil, fmt.Errorf("error parsing crime values from row"+
			": %s", err.Error())
	}

	// Parse date range
	startRange, endRange, err := date.NewRangeFromStr(rangeStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing date range into "+
			"time.Times, err: %s", err.Error())
	}
	crime.DateOccurredStart = startRange
	crime.DateOccurredEnd = endRange

	// Success
	return crime, nil
}

func (c Crime) String() string {
	return fmt.Sprintf("ReportID: %d\n"+
		"Page: %d\n"+
		"Reported: %s\n"+
		"Occurred Start: %s\n"+
		"Occurred End: %s\n"+
		"ID: %d-%d\n"+
		"GeoLocID: %d\n"+
		"Incidents: %s\n"+
		"Description: %s\n"+
		"Remediation: %s\n"+
		"Parse Errors: %s",
		c.ReportID,
		c.Page,
		c.DateReported,
		c.DateOccurredStart,
		c.DateOccurredEnd,
		c.ReportSuperID,
		c.ReportSubID,
		c.GeoLocID,
		strings.Join(c.Incidents, ","),
		strings.Join(c.Descriptions, ","),
		c.Remediation,
		strings.Join(StringParseErrors(c.ParseErrors), ", "))
}

// Query finds a model with matching attributes in the db and sets the Crime.ID
// field if found. Additionally an error is returned. Which will be
// sql.ErrNoRows if a matching model is not found. Or nil on success.
func (c *Crime) Query() error {
	// Get db
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error getting db instance: %s", err.Error())
	}

	// Query
	row := db.QueryRow("SELECT id FROM crimes WHERE report_id=$1 AND "+
		"page=$2 AND date_reported=$3 AND date_occurred=tstzrange($4,"+
		" $5, '()') AND report_super_id=$6 AND report_sub_id=$7 AND "+
		"incidents=$8 AND descriptions=$9 AND "+
		"remediation=$10", c.ReportID, c.Page, c.DateReported,
		c.DateOccurredStart, c.DateOccurredEnd, c.ReportSuperID,
		c.ReportSubID, c.Incidents, c.Descriptions,
		c.Remediation)

	// Get ID
	err = row.Scan(&c.ID)

	// Check if no row found
	if err == sql.ErrNoRows {
		// Just return sql error so we can ID
		return err
	} else if err != nil {
		return fmt.Errorf("error querying for crime model: %s",
			err.Error())
	}

	return nil
}

// Insert adds the model to the database and sets the Crime.ID field to the
// newly inserted models ID. Additionally an error is returned if one occurs,
// or nil on success.
func (c *Crime) Insert() error {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error creating db instance: %s",
			err.Error())
	}

	// Insert
	row := db.QueryRow("INSERT INTO crimes (report_id, page, date_reported, "+
		"date_occurred, report_super_id, report_sub_id, geo_loc_id, "+
		"incidents, descriptions, remediation) VALUES ($1, $2, $3, "+
		"tstzrange($4, $5, '()'), $6, $7, $8, $9, $10, $11) RETURNING id",
		c.ReportID, c.Page, c.DateReported, c.DateOccurredStart,
		c.DateOccurredEnd, c.ReportSuperID, c.ReportSubID, c.GeoLocID,
		c.Incidents, c.Descriptions, c.Remediation)

	// Get ID
	err = row.Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("error inserting into db: %s",
			err.Error())
	}

	return nil
}

// InsertIfNew saves the current Crime model if it does not exist in the db.
// Returns an error if one occurs, or nil on success.
func (c *Crime) InsertIfNew() error {
	// Query
	err := c.Query()

	// Check if doesn't exist
	if err == sql.ErrNoRows {
		// Insert
		err := c.Insert()
		if err != nil {
			return fmt.Errorf("error inserting non existing model: %s",
				err.Error())
		}
	} else if err != nil {
		// General error
		return fmt.Errorf("error determining if crime is new: %s",
			err.Error())
	}

	// Success
	return nil
}

// QueryAllCrimes retrieves the specified number of Crime models from the
// database. Ordered by the field specified in the orderBy argument. Must be
// one of 'date_reported' or 'date_occurred'. An array of Crimes are returned,
// along with an error. Which is nil on success.
//
// Retrieves all crime columns.
func QueryAllCrimes(offset uint, limit uint, orderBy OrderByType) ([]*Crime, error) {
	crimes := []*Crime{}

	// Check orderBy var
	if orderBy == OrderByErr {
		return crimes, fmt.Errorf("invalid orderBy value: %s", orderBy)
	}

	// Get db
	db, err := dstore.NewDB()
	if err != nil {
		return crimes, fmt.Errorf("error retrieving database instance"+
			": %s", err.Error())
	}

	// Query
	rows, err := db.Query("SELECT id, report_id, page, date_reported, "+
		"date_occurred, report_super_id, report_sub_id, "+
		"geo_loc_id, incidents, descriptions, remediation "+
		"FROM crimes ORDER BY $1 DESC OFFSET $2 LIMIT $3",
		orderBy, offset, limit)

	if err != nil {
		return crimes, fmt.Errorf("error querying database for crimes"+
			": %s", err.Error())
	}

	// Parse
	for rows.Next() {
		crime, err := NewCrime(rows)
		if err != nil {
			return crimes, fmt.Errorf("error parsing crime row: %s",
				err.Error())
		}

		crimes = append(crimes, crime)
	}

	// Close query
	if err = rows.Close(); err != nil {
		return crimes, fmt.Errorf("error closing crimes query: %s",
			err.Error())
	}

	// Success
	return crimes, nil
}
