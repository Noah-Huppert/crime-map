package models

import (
	"fmt"
	dbp "github.com/Noah-Huppert/crime-map/db"
)

// Migrate will attempt to create all tables defined by models. And return an
// error if one occurs, nil otherwise.
func Migrate() error {
	// Make db instance
	db, err := dbp.NewDB()
	if err != nil {
		return fmt.Errorf("error making db instance: %s", err.Error())
	}

	// Migrate
	err = db.AutoMigrate(&Crime{}).Error
	if err != nil {
		return fmt.Errorf("error migrating db: %s", err.Error())
	}

	return nil
}
