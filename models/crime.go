package models

import (
	"database/sql"
	"fmt"
	"github.com/Noah-Huppert/crime-map/dstore"
	"github.com/lib/pq"
	"strings"
	"time"
)

// Crime structs hold information about criminal activity reported by Clery
// act reports
type Crime struct {
	// ID is a unique identifier
	ID int

	// University indicates which university the crime was reported by
	University string

	// DateReported records when the criminal activity was disclosed to the
	// police
	DateReported time.Time

	// DateOccurredStart records when the criminal activity started taking
	// place
	DateOccurredStart time.Time

	// DateOccurredEnd records when the criminal activity stopped taking
	// place
	DateOccurredEnd time.Time

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
	ParseErrors []ParseError
}

func (c Crime) String() string {
	return fmt.Sprintf("University: %s\n"+
		"Reported: %s\n"+
		"Occurred Start: %s\n"+
		"Occurred End: %s\n"+
		"ID: %d-%d\n"+
		"GeoLocID: %d\n"+
		"Incidents: %s\n"+
		"Description: %s\n"+
		"Remediation: %s\n"+
		"Parse Errors: %s",
		c.University,
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
	row := db.QueryRow("SELECT id FROM crimes WHERE university=$1 AND "+
		"date_reported=$2 AND date_occurred=tstzrange($3, $4, '()') "+
		"AND report_super_id=$5 AND report_sub_id=$6 AND "+
		"incidents=$7 AND descriptions=$8 AND "+
		"remediation=$9", c.University, c.DateReported,
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
	row := db.QueryRow("INSERT INTO crimes (university, date_reported, "+
		"date_occurred, report_super_id, report_sub_id, geo_loc_id, "+
		"incidents, descriptions, remediation) VALUES ($1, $2, "+
		"tstzrange($3, $4, '()'), $5, $6, $7, $8, $9, $10) RETURNING id",
		c.University, c.DateReported, c.DateOccurredStart,
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

// SaveIfNew saves the current Crime model if it does not exist in the db.
// Returns an error if one occurs, or nil on success.
func (c *Crime) SaveIfNew() error {
	// Query
	err := c.Query()
	if (err != nil) && (err != sql.ErrNoRows) {
		return fmt.Errorf("error determining if crime is new: %s",
			err.Error())
	}

	// Check if doesn't exist
	if err == sql.ErrNoRows {
		// Insert
		err := c.Insert()
		if err != nil {
			return fmt.Errorf("error inserting non existing model: %s",
				err.Error())
		}
	}

	// Success
	return nil
}
