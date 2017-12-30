package models

import (
	"fmt"
	"github.com/Noah-Huppert/crime-map/dstore"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

// Migrate will attempt to create all tables defined by models. And return an
// error if one occurs, nil otherwise.
func Migrate() error {
	// Make db instance
	db, err := dstore.NewDB()
	if err != nil {
		return fmt.Errorf("error making db instance: %s", err.Error())
	}

	// Make db driver for migration
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error making db driver: %s", err.Error())
	}

	// Create migrator
	migrator, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("error making migrator instance: %s",
			err.Error())
	}

	// Run
	if err = migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running migrations: %s", err.Error())
	}

	return nil
}
