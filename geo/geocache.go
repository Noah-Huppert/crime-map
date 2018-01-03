package geo

import (
	"database/sql"
	"fmt"

	"github.com/Noah-Huppert/crime-map/models"
)

// GeoCache caches GeoLoc models retrieved from the database
type GeoCache struct {
	// locs holds all GeoLoc models retrieved from the database
	locs map[string]*models.GeoLoc
}

// NewGeoCache constructs a new GeoCache object
func NewGeoCache() *GeoCache {
	return &GeoCache{
		locs: make(map[string]*models.GeoLoc),
	}
}

// Get retrieves a GeoCache model with the provided raw value. This model will
// be populated with the raw and ID field only. An error is returned if one
// occurs, or nil on success.
func (c GeoCache) Get(raw string) (*models.GeoLoc, error) {
	// Check cached in locs var
	if val, ok := c.locs[raw]; ok {
		return val, nil
	}

	// If not cached
	// Make GeoLoc model with provided raw field
	loc := models.NewGeoLoc(raw)

	// Query
	err := loc.Query()

	// Check if not found
	if err == sql.ErrNoRows {
		// Return err we can identify there are no rows. AND return
		// the GeoLoc model we contructed so we can save it into the db
		return loc, err
	} else if err != nil {
		return nil, fmt.Errorf("error querying database for GeoLoc"+
			", raw: \"%s\", err: %s", raw, err.Error())
	}

	// If found, cache in locs var
	c.locs[raw] = loc

	// Success
	return loc, nil
}

// InsertIfNew inserts the GeoLoc model into the database if it does not exist.
// The ID of the model in the database will be set in the GeoLoc.ID field. An
// error is returned if one occurs, nil on success.
func (c GeoCache) InsertIfNew(raw string) (*models.GeoLoc, error) {
	// Query
	loc, err := c.Get(raw)

	// Check if model doesn't exist
	if err == sql.ErrNoRows {
		// Insert
		if err = loc.Insert(); err != nil {
			return nil, fmt.Errorf("error inserting non-existent GeoLoc"+
				" model: %s", err.Error())
		}
	} else if err != nil {
		// General error
		return nil, fmt.Errorf("error querying for GeoLoc model: %s",
			err.Error())
	}

	// Success
	return loc, nil
}
