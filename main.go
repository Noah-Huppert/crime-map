package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Noah-Huppert/crime-map/geo"
	"github.com/Noah-Huppert/crime-map/http"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/parsers"
)

const file = "data/2017-10-12.pdf"

func main() {
	// Make context to control running of async jobs
	ctx := context.Background()

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

	// Save crimes
	fmt.Println("saving crimes")
	for _, crime := range crimes {
		if err = crime.InsertIfNew(); err != nil {
			fmt.Printf("error saving crime, crime: %s, err: %s\n",
				crime,
				err.Error())
			os.Exit(1)
			return
		}

		// Save any parse errors
		for _, pErr := range crime.ParseErrors {
			// Set Crime FK
			pErr.CrimeID = crime.ID

			// Save
			if err = pErr.InsertIfNew(); err != nil {
				fmt.Printf("error saving crime parse error, "+
					"crime: %s, parse err: %s, err: %s",
					crime, pErr, err.Error())
				os.Exit(1)
				return
			}
		}
	}

	// Locate
	fmt.Println("querying for unlocated GeoLoc models")
	locater := geo.NewLocater()
	unlocated, err := models.QueryUnlocatedGeoLocs()

	if err != nil {
		fmt.Printf("error querying for unlocated GeoLocs: %s\n",
			err.Error())
		os.Exit(1)
		return
	}

	fmt.Println("locating unlocated GeoLoc models")
	errs := make(chan error)
	locs := make(chan *models.GeoLoc)
	currentlyLocating := len(unlocated)

	// Kick off locator jobs
	for _, loc := range unlocated {
		// Locate
		locater.LocateAsync(ctx, errs, locs, loc)
	}

	for currentlyLocating > 0 {
		select {
		case err = <-errs:
			// If error
			currentlyLocating -= 1
			fmt.Printf("error locating model: %s", err.Error())
			os.Exit(1)
			return
		case loc := <-locs:
			// If success
			currentlyLocating -= 1

			if err = loc.Update(); err != nil {
				fmt.Printf("error updating GeoLoc model, loc: %s, "+
					"err: %s\n",
					loc, err.Error())
				os.Exit(1)
				return
			}

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
