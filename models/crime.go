package models

import (
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
	ID uint

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

	// ReportID is the second portion of the police report ID associated
	// with the reported crime
	// TODO: Rename to report sub id
	ReportID uint

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
		"Remediation: %s",
		c.DateReported,
		c.DateOccurredStart,
		c.DateOccurredEnd,
		c.ReportSuperID,
		c.ReportID,
		c.Location,
		c.GeoLocID,
		strings.Join(c.Incidents, ","),
		strings.Join(c.Descriptions, ","),
		c.Remediation)
}

// Query finds a model with matching attributes in the db and returns the db
// rows object. Additionally an error is returned if one occurs, or nil on
// success.
func (c Crime) Query() (*db.Rows, error) {
	// Get db
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error getting db instance: %s", err.Error())
	}

	// Query
	rows, err := db.Query("SELECT id FROM crimes WHERE date_reported=$1 " +
		"AND date_occurred=[$2, $3] AND " +
		"report_super_id=$4 AND report_sub_id=$5 AND " +
		"location=$6 AND geo_loc_id=$7 AND " +
		"incidents=$8 AND descriptions=$9 AND " +
		"remediation=$10")

	if err != nil {
		return nil, fmt.Errorf("error querying database: %s",
			err.Error())
	}

	return rows, nil
}

// Insert adds the model to the database
func (c Crime) Insert() (*db.Result, error) {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return nil, fmt.Errorf("error creating db instance: %s", 
			err.Error())
	}

	// Insert
	res, err := db.Exec("INSERT INTO crimes (id, date_reported, date_occurred)
	// TODO: Finish insert
}

// SaveIfUnique saves the current Crime model if it does not exist in the db.
// Returns an error if one occurs, or nil on success.
func (c Crime) SaveIfUnique() error {
	// Query
	rows, err := c.Query()

	// Check if exists
	if len(rows) == 0 {
		// If doesn't, insert
		res, err 
	}
}
