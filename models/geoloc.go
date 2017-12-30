package models

import (
	"database/sql"
	"fmt"

	"github.com/Noah-Huppert/crime-map/dstore"
)

const (
	// StatusOk indicates that a GAPI request was successful
	StatusOk string = "OK"

	// StatusZero indicates that the GAPI request didn't return any results
	StatusZero string = "ZERO_RESULTS"

	// StatusLimit indicates that we have gone over the GAPI query limit
	StatusLimit string = "OVER_QUERY_LIMIT"

	// StatusDenied indicates that the GAPI request was denied
	StatusDenied string = "REQUEST_DENIED"

	// StatusInvalid indicates that the GAPI request was invalid
	StatusInvalid string = "INVALID_REQUEST"

	// StatusErr indicates that an unknown error occurred during the GAPI
	// request
	StatusErr string = "UNKNOWN_ERROR"
)

const (
	// AccuracyPerfect indicates that the location provided by the GAPI
	// is exact
	AccuracyPerfect string = "ROOFTOP"

	// AccuracyBetween indicates that the location provided by the GAPI
	// is between two addresses
	AccuracyBetween string = "RANGE_INTERPOLATED"

	// AccuracyCenter indicates that the location is in the middle of an
	// region. Such as a block
	AccuracyCenter string = "GEOMETRIC_CENTER"

	// AccuracyApprox indicates that the location is not exact
	AccuracyApprox string = "APPROXIMATE"
)

// GeoLoc holds information about the geographical location of a crime "location"
// field.
//
// GeoLoc are resolved by using the Google Geocoding API to transform location
// strings into lat long coords. GeoLoc also holds some additional accuracy
// information, as not all locations can be resolved exactly.
type GeoLoc struct {
	// ID is the unique identifier
	ID int

	// Located indicates if the raw location has been geocoded using the
	// GAPI
	Located bool

	// Lat is the latitude of the location
	Lat float32

	// Long is the longitude of the location
	Long float32

	// PostalAddr holds the formatted postal address of the location
	PostalAddr string

	// Accuracy indicates how close to the provided location the lat long
	// are
	Accuracy string

	// Partial indicates if the match is only a partial
	Partial bool

	// BoundsProvided indicates whether any location bounds were provided
	BoundsProvided bool

	// BoundsID holds the GeoBounds ID representing the area the location
	// covers
	BoundsID uint

	// GAPIPlaceID holds the GAPI location ID, used to retrieve additional
	// information about a location using the GAPI
	GAPIPlaceID string

	// GAPIStatus holds the status of the GAPI Geocoding request which was
	// sent to locate the location
	GAPIStatus string

	// Raw holds the text present on the crime report which the GeoLoc
	// attempts to locate
	Raw string
}

// NewGeoLoc returns a new GeoLoc instance with the provided raw text
func NewGeoLoc(raw string) *GeoLoc {
	return &GeoLoc{
		Located: false,
		Raw:     raw,
	}
}

func (l GeoLoc) String() string {
	return fmt.Sprintf("ID: %d\n"+
		"Located: %t\n"+
		"Lat: %f\n"+
		"Long: %f\n"+
		"PostalAddr: %s\n"+
		"Accuracy: %s\n"+
		"Partial: %t\n"+
		"BoundsProvided: %t\n"+
		"BoundsID: %d\n"+
		"GAPIPlaceID: %s\n"+
		"Raw: %s",
		l.ID, l.Located, l.Lat, l.Long, l.PostalAddr,
		l.Accuracy, l.Partial, l.BoundsProvided, l.BoundsID,
		l.GAPIPlaceID, l.Raw)
}

// Query attempts to find a GeoLoc model in the db with the same raw field
// value. If a model is found, the GeoLoc.ID field is set. Additionally an
// error is returned if one occurs. sql.ErrNoRows is returned if no GeoLocs
// were found. Or nil on success.
func (l *GeoLoc) Query() error {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving db instance: %s",
			err.Error())
	}

	// Query
	row := db.QueryRow("SELECT id FROM geo_locs WHERE raw = $1", l.Raw)

	// Get ID
	err = row.Scan(&l.ID)
	fmt.Printf("l.ID: %d\n", l.ID)

	// Check if row found
	if err == sql.ErrNoRows {
		// If not, return so we can identify
		return err
	} else if err != nil {
		return fmt.Errorf("error reading GeoLoc ID from row: %s",
			err.Error())
	}

	// Success
	return nil
}

// Insert adds a GeoLoc model to the database. An error is returned if one
// occurs, or nil on success.
func (l *GeoLoc) Insert() error {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving DB instance: %s",
			err.Error())
	}

	// Insert
	var row *sql.Row

	// Check if GeoLoc has been parsed
	if l.Located {
		// If so, save all fields
		row = db.QueryRow("INSERT INTO geo_locs (located, lat, long,"+
			" postal_addr, accuracy, partial, bounds_provided, "+
			"bounds_id, gapi_place_id, raw) VALUES ($1, $2, $3, $4"+
			", $5, $6, $7, $8, $9, $10) RETURNING id",
			l.Located, l.Lat, l.Long, l.PostalAddr, l.Accuracy,
			l.Partial, l.BoundsProvided, l.BoundsID, l.GAPIPlaceID,
			l.Raw)
	} else {
		// If not, only save a couple, and leave rest null
		row = db.QueryRow("INSERT INTO geo_locs (located, raw) VALUES"+
			" ($1, $2) RETURNING id",
			l.Located, l.Raw)
	}

	// Get inserted row ID
	err = row.Scan(&l.ID)
	if err != nil {
		return fmt.Errorf("error inserting row, Located: %t, err: %s",
			l.Located, err.Error())
	}

	return nil
}

/*
// InsertIfNew adds the GeoLoc model to the database if a model with its fields
// does not exist. An error is returned if one occurs, or nil on success.
//
// Uses the provided GeoCache to query for the existence of a model.
func (l GeoLoc) InsertIfNew(geoCache *GeoCache) error {
	l, err := geoCache

	// Check if not found
	if err == sql.ErrNoRows {
		// Insert
		err = l.Insert()

		if err != nil {
			return fmt.Errorf("error inserting geo loc model: %s",
				err.Error())
		}
	}

	// Success
	return nil
}
*/
