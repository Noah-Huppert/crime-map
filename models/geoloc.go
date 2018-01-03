package models

import (
	"database/sql"
	"fmt"

	"github.com/Noah-Huppert/crime-map/dstore"
)

type GeoLocAccuracy string

const (
	// AccuracyPerfect indicates that the location provided by the GAPI
	// is exact
	AccuracyPerfect GeoLocAccuracy = "ROOFTOP"

	// AccuracyBetween indicates that the location provided by the GAPI
	// is between two addresses
	AccuracyBetween GeoLocAccuracy = "RANGE_INTERPOLATED"

	// AccuracyCenter indicates that the location is in the middle of an
	// region. Such as a block
	AccuracyCenter GeoLocAccuracy = "GEOMETRIC_CENTER"

	// AccuracyApprox indicates that the location is not exact
	AccuracyApprox GeoLocAccuracy = "APPROXIMATE"

	// AccuracyErr indicates that an invalid string value was provided
	// when creating a GeoLocAccuracy
	AccuracyErr GeoLocAccuracy = "ERR"
)

func NewGeoLocAccuracy(str string) (GeoLocAccuracy, error) {
	if str == string(AccuracyPerfect) {
		return AccuracyPerfect, nil
	} else if str == string(AccuracyBetween) {
		return AccuracyBetween, nil
	} else if str == string(AccuracyCenter) {
		return AccuracyCenter, nil
	} else if str == string(AccuracyApprox) {
		return AccuracyApprox, nil
	} else {
		return AccuracyErr, fmt.Errorf("unknown GeoLocAccuracy string "+
			", str: %s", str)
	}
}

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

	// GAPISuccess indicates if the GAPI locate request succeeded
	GAPISuccess bool

	// Lat is the latitude of the location
	Lat float64

	// Long is the longitude of the location
	Long float64

	// PostalAddr holds the formatted postal address of the location
	PostalAddr string

	// Accuracy indicates how close to the provided location the lat long
	// are
	Accuracy GeoLocAccuracy

	// BoundsProvided indicates whether any location bounds were provided
	BoundsProvided bool

	// BoundsID holds the GeoBounds ID which specifies the location of the
	// crime
	BoundsID sql.NullInt64

	// ViewportBoundsID holds the GeoBounds ID which specifies the
	// recommended viewport for looking at the crime location
	ViewportBoundsID int

	// GAPIPlaceID holds the GAPI location ID, used to retrieve additional
	// information about a location using the GAPI
	GAPIPlaceID string

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

// NewUnlocatedGeoLoc creates a new GeoLoc instance from the currently selected
// result set in the provided sql.Rows object. This row should select the id
// and raw fields, in that order. An error will be returned if one occurs, nil
// on success.
func NewUnlocatedGeoLoc(row *sql.Rows) (*GeoLoc, error) {
	loc := NewGeoLoc("")

	// Parse
	if err := row.Scan(&loc.ID, &loc.Raw); err != nil {
		return nil, fmt.Errorf("error reading field values from row: %s",
			err.Error())
	}

	// Success
	return loc, nil
}

func (l GeoLoc) String() string {
	return fmt.Sprintf("ID: %d\n"+
		"Located: %t\n"+
		"GAPISuccess: %t\n"+
		"Lat: %f\n"+
		"Long: %f\n"+
		"PostalAddr: %s\n"+
		"Accuracy: %s\n"+
		"BoundsProvided: %t\n"+
		"BoundsID: %d\n"+
		"ViewportBoundsID: %d\n"+
		"GAPIPlaceID: %s\n"+
		"Raw: %s",
		l.ID, l.Located, l.GAPISuccess, l.Lat, l.Long, l.PostalAddr,
		l.Accuracy, l.BoundsProvided, l.BoundsID, l.ViewportBoundsID,
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

// Update sets an existing GeoLoc model's fields to new values. Only updates the
// located and raw fields if located == false. Updates all fields if
// located == true.
//
// It relies on the raw field to specify exactly which row to update. The row
// column has a unique constraint, so this is sufficient.
//
// An error is returned if one occurs, or nil on success.
func (l GeoLoc) Update() error {
	// Get database instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Update
	var row *sql.Row

	// If not located
	if !l.Located {
		row = db.QueryRow("UPDATE geo_locs SET located = $1, raw = "+
			"$2 WHERE raw = $2 RETURNING id", l.Located, l.Raw)
	} else {
		// Check accuracy value
		if l.Accuracy == AccuracyErr {
			return fmt.Errorf("invalid accuracy value: %s",
				l.Accuracy)
		}
		// If located
		row = db.QueryRow("UPDATE geo_locs SET located = $1, "+
			"gapi_success = $2, lat = $3, long = $4, "+
			"postal_addr = $5, accuracy = $6, bounds_provided = $7,"+
			"bounds_id = $8, viewport_bounds_id = $9, "+
			"gapi_place_id = $10, raw = $11 WHERE raw = $11 "+
			"RETURNING id",
			l.Located, l.GAPISuccess, l.Lat, l.Long, l.PostalAddr,
			l.Accuracy, l.BoundsProvided, l.BoundsID,
			l.ViewportBoundsID, l.GAPIPlaceID, l.Raw)
	}

	// Set ID
	err = row.Scan(&l.ID)

	// If doesn't exist
	if err == sql.ErrNoRows {
		// Return error so we can identify
		return err
	} else if err != nil {
		// Other error
		return fmt.Errorf("error updating GeoLoc, located: %t, err: %s",
			l.Located, err.Error())
	}

	// Success
	return nil
}

// QueryUnlocatedGeoLocs finds all GeoLoc models which have not been located on
// a map. Additionally an error is returned if one occurs, or nil on success.
func QueryUnlocatedGeoLocs() ([]*GeoLoc, error) {
	locs := []*GeoLoc{}

	// Get db
	db, err := dstore.NewDB()
	if err != nil {
		return locs, fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Query
	rows, err := db.Query("SELECT id, raw FROM geo_locs WHERE located = " +
		"false")

	// Check if no results
	if err == sql.ErrNoRows {
		// If not results, return raw error so we can identify
		return locs, err
	} else if err != nil {
		// Other error
		return locs, fmt.Errorf("error querying for unlocated GeoLocs"+
			": %s", err.Error())
	}

	// Parse rows into GeoLocs
	for rows.Next() {
		// Parse
		loc, err := NewUnlocatedGeoLoc(rows)
		if err != nil {
			return locs, fmt.Errorf("error creating unlocated "+
				"GeoLoc from row: %s", err.Error())
		}

		// Add to list
		locs = append(locs, loc)
	}

	// Close
	if err = rows.Close(); err != nil {
		return locs, fmt.Errorf("error closing query: %s",
			err.Error())
	}

	// Success
	return locs, nil
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
		// Check accuracy value
		if l.Accuracy == AccuracyErr {
			return fmt.Errorf("invalid accuracy value: %s",
				l.Accuracy)
		}

		// If so, save all fields
		row = db.QueryRow("INSERT INTO geo_locs (located, gapi_success"+
			", lat, long, postal_addr, accuracy, bounds_provided, "+
			"bounds_id, viewport_bounds_id, gapi_place_id, raw) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) "+
			"RETURNING id",
			l.Located, l.GAPISuccess, l.Lat, l.Long, l.PostalAddr,
			l.Accuracy, l.BoundsProvided, l.BoundsID,
			l.ViewportBoundsID, l.GAPIPlaceID, l.Raw)
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
