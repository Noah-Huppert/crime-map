package dstore

import (
	"fmt"
	"github.com/Noah-Huppert/crime-map/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// instance holds the database instance if already created
var instance *gorm.DB

// NewDB creates a new connected DB instance and returns it. Along with an
// error if one occurs. Or nil on success.
func NewDB() (*gorm.DB, error) {
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
	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
			c.DB.Host, c.DB.User, c.DB.Password, c.DB.Name,
			c.DB.SSLMode))

	if err != nil {
		return nil, fmt.Errorf("error connecting to db: %s", err.Error())
	}

	instance = db

	return instance, nil
}
