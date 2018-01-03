package models

import (
	"database/sql"
	"fmt"
	"googlemaps.github.io/maps"

	"github.com/Noah-Huppert/crime-map/dstore"
)

// GeoBound indicates a square area on a map
type GeoBound struct {
	// ID is the unique identifier
	ID int

	// NeLat holds the latitude which the northeast corner of the map
	// viewport should be located at
	NeLat float64

	// NeLong holds the longitude which the northeast corner of the map
	// viewport should be located at
	NeLong float64

	// SwLat holds the latitude which the southwest corner of the map
	// viewport should be located at
	SwLat float64

	// SwLong holds the longitude which the southwest corner of the map
	// viewport should be located at
	SwLong float64
}

// GeoBoundFromMapsBound creates a new GeoBound instance from a Google Maps API
// maps.LatLngBounds structure
func GeoBoundFromMapsBound(bounds maps.LatLngBounds) *GeoBound {
	return &GeoBound{
		NeLat:  bounds.NorthEast.Lat,
		NeLong: bounds.NorthEast.Lng,
		SwLat:  bounds.SouthWest.Lat,
		SwLong: bounds.SouthWest.Lng,
	}
}

// Query attempts to locate a GeoBound with the same Ne and Sw Lat Long values
// in the database. The GeoBound.ID field will be populated with the model's
// ID in the database. An error will be returned if one occurs, or nil on
// success.
func (b *GeoBound) Query() error {
	// Get database instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Query
	row := db.QueryRow("SELECT id FROM geo_bounds WHERE ne_lat = $1 AND "+
		"ne_long = $2 AND sw_lat = $3 AND sw_long = $4",
		b.NeLat, b.NeLong, b.SwLat, b.SwLong)

	// Get ID
	err = row.Scan(&b.ID)
	// If no rounds found
	if err == sql.ErrNoRows {
		// Just return error so we can identify
		return err
	} else if err != nil {
		return fmt.Errorf("error querying db for GeoBound model: %s",
			err.Error())
	}

	return nil
}

// Insert adds a GeoBound model to the database. An error is returned if one
// occurs, nil on success.
func (b *GeoBound) Insert() error {
	// Get db instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error retrieving database instance: %s",
			err.Error())
	}

	// Insert
	row := db.QueryRow("INSERT INTO geo_bounds (ne_lat, ne_long, sw_lat, "+
		"sw_long) VALUES ($1, $2, $3, $4) RETURNING id",
		b.NeLat, b.NeLong, b.SwLat, b.SwLong)

	// Get new ID
	err = row.Scan(&b.ID)
	if err != nil {
		return fmt.Errorf("error inserting GeoBound into db: %s",
			err.Error())
	}

	// Success
	return nil
}

// InsertIfNew attempts to find a GeoBound model with the same values in the
// database. If none is found, the model is added to the database. In both
// cases the GeoBound.ID field is set to that of the found/inserted row in the
// db. An error is returned if one occurs, or nil on success.
func (b *GeoBound) InsertIfNew() error {
	// Query
	err := b.Query()

	// If doesn't exist yet
	if err == sql.ErrNoRows {
		// Insert
		if err = b.Insert(); err != nil {
			return fmt.Errorf("error inserting non existing "+
				"GeoBound: %s", err.Error())
		}
	} else if err != nil {
		// General error
		return fmt.Errorf("error querying for GeoBound: %s",
			err.Error())
	}

	// Success
	return nil
}
