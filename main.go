package main

import (
	"fmt"
	"os"

	"github.com/Noah-Huppert/crime-map/db"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/parsers"
)

const file = "data/2017-10-12.pdf"

func main() {
	// Migrate db
	err := models.Migrate()
	if err != nil {
		fmt.Printf("error migrating db: %s", err.Error())
		os.Exit(1)
	}

	// Connect to db
	db, err := db.NewDB()
	if err != nil {
		fmt.Printf("error creating db instance: %s", err.Error())
		os.Exit(1)
	}

	// Parse crimes
	r := parsers.NewReport(file)

	crimes, err := r.Parse()
	if err != nil {
		fmt.Printf("error parsing report: %s", err.Error())
		os.Exit(1)
	}

	// Print crimes
	for i, crime := range crimes {
		fmt.Printf("\n%d\n====\n%s\n", i+1, crime)

		if err = db.Save(&crime).Error; err != nil {
			fmt.Printf("error saving crime: %s\n", err.Error())
			os.Exit(1)
		}
	}
}
