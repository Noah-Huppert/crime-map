package main

import (
	"fmt"
	"os"

	"github.com/Noah-Huppert/crime-map/geo"
	"github.com/Noah-Huppert/crime-map/http"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/parsers"
)

const file = "data/2017-10-12.pdf"

func main() {
	// Migrate db
	fmt.Println("migrating db")
	err := models.Migrate()
	if err != nil {
		fmt.Printf("error migrating db: %s\n", err.Error())
		os.Exit(1)
		return
	}

	// Make geocache
	geoCache := geo.NewGeoCache()

	// Parse crimes
	fmt.Println("parsing report")
	r := parsers.NewReport(file, geoCache)

	crimes, err := r.Parse()
	if err != nil {
		fmt.Printf("error parsing report: %s\n", err.Error())
		os.Exit(1)
		return
	}

	// Print crimes
	for _, crime := range crimes {
		fmt.Printf("%s\n", crime)
		if err = crime.SaveIfNew(); err != nil {
			fmt.Printf("error saving crime: %s\n", err.Error())
			os.Exit(1)
			return
		}
	}

	// Start http server
	err = http.Serve()
	if err != nil {
		fmt.Printf("error starting http server: %s\n", err.Error())
		os.Exit(1)
		return
	}
}
