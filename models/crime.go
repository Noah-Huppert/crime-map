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

	// TODO: Add "University" field

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

	// Location is the place where the criminal activity occurred
	Location string

	// GeoLocID is the unique ID of the Geo entry which holds the geographically
	// encoded location in lat long form
	GeoLocID uint `gorm:"ForeignKey:GeoLocID"`

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
	return fmt.Sprintf("Reported: %s\n"+
		"Occurred Start: %s\n"+
		"Occurred End: %s\n"+
		"ID: %d-%d\n"+
		"Location: %s\n"+
		"GeoLocID: %d\n"+
		"Incidents: %s\n"+
		"Description: %s\n"+
		"Remediation: %s\n"+
		"Parse Errors: %s",
		c.DateReported,
		c.DateOccurredStart,
		c.DateOccurredEnd,
		c.ReportSuperID,
		c.ReportSubID,
		c.Location,
		c.GeoLocID,
		strings.Join(c.Incidents, ","),
		strings.Join(c.Descriptions, ","),
		c.Remediation,
		strings.Join(StringParseErrors(c.ParseErrors), ", "))
}

// Query finds a model with matching attributes in the db and returns the db
// rows object. Which must be closed. Additionally an error is returned if one
// occurs, or nil on success.
func (c Crime) Query() (*sql.Rows, error) {
	// Get db
	db, err := dstore.NewDB()
	if err != nil {
		return nil, fmt.Errorf("error getting db instance: %s", err.Error())
	}

	// Query
	rows, err := db.Query("SELECT id FROM crimes WHERE date_reported=$1 "+
		"AND date_occurred=tstzrange($2, $3, '()') AND "+
		"report_super_id=$4 AND report_sub_id=$5 AND "+
		"location=$6 AND geo_loc_id=$7 AND "+
		"incidents=$8 AND descriptions=$9 AND "+
		"remediation=$10", c.DateReported, c.DateOccurredStart,
		c.DateOccurredEnd, c.ReportSuperID, c.ReportSubID, c.Location,
		c.GeoLocID, c.Incidents, c.Descriptions, c.Remediation)

	if err != nil {
		return nil, fmt.Errorf("error querying database: %s",
			err.Error())
	}

	return rows, nil
}

// Insert adds the model to the database. Returns a db rows object for the
// insert containing the new model id. This rows object must be closed.
// Additionally an error is returned, nil on success.
func (c Crime) Insert() (*sql.Rows, error) {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return nil, fmt.Errorf("error creating db instance: %s",
			err.Error())
	}

	// Insert
	rows, err := db.Query("INSERT INTO crimes (date_reported, "+
		"date_occurred, report_super_id, report_sub_id, location, "+
		"geo_loc_id, incidents, descriptions, remediation) VALUES "+
		"($1, tstzrange($2, $3, '()'), $4, $5, $6, NULL, $7, $8, "+
		"$9) RETURNING id",
		c.DateReported, c.DateOccurredStart, c.DateOccurredEnd,
		c.ReportSuperID, c.ReportSubID, c.Location,
		c.Incidents, c.Descriptions, c.Remediation)

	if err != nil {
		return nil, fmt.Errorf("error inserting into db: %s",
			err.Error())
	}

	return rows, nil
}

// SaveIfNew saves the current Crime model if it does not exist in the db.
// Returns an error if one occurs, or nil on success.
func (c Crime) SaveIfNew() error {
	// Query
	rows, err := c.Query()
	if err != nil {
		return fmt.Errorf("error determining if crime is new: %s",
			err.Error())
	}

	// Check if exists
	if !rows.Next() {
		// If doesn't exist, close query
		if err = rows.Close(); err != nil {
			return fmt.Errorf("error closing query when model "+
				"doesn't exist: %s", err.Error())
		}

		// Insert
		res, err := c.Insert()
		if err != nil {
			return fmt.Errorf("error inserting non existing model: %s",
				err.Error())
		}

		// Get new id
		if !res.Next() {
			return fmt.Errorf("error retrieving crime model id, " +
				"query didn't return any rows")
		}

		if err = res.Scan(&c.ID); err != nil {
			return fmt.Errorf("error retrieving model id from "+
				"insert rows: %s",
				err.Error())
		}

		// Close
		if err = res.Close(); err != nil {
			return fmt.Errorf("error closing insert query: %s",
				err.Error())
		}
	} else {
		// Already exists, get id
		if err = rows.Scan(&c.ID); err != nil {
			return fmt.Errorf("error retrieving model id from query: %s",
				err.Error())
		}
	}

	// Close query
	if err = rows.Close(); err != nil {
		return fmt.Errorf("error closing query when model exists: %s",
			err.Error())
	}

	// Success
	return nil
}
