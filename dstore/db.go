package dstore

import (
	"database/sql"
	"fmt"
	"github.com/Noah-Huppert/crime-map/config"
	_ "github.com/lib/pq"
)

// instance holds the database instance if already created
var instance *sql.DB

// NewDB creates a new connected DB instance and returns it. Along with an
// error if one occurs. Or nil on success.
func NewDB() (*sql.DB, error) {
	// Check if exists
	if instance != nil {
		return instance, nil
	}

	// Get db config
	c, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading configuration: %s",
			err.Error())
	}

	// Connect
	db, err := sql.Open("postgres", c.DB.ConnString)

	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %s", err.Error())
	}

	// Verify
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error testing db connectivity: %s",
			err.Error())
	}

	instance = db

	return instance, nil
}
