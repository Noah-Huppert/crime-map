package geo

import (
	"context"
	"fmt"

	"github.com/Noah-Huppert/crime-map/models"
)

// LocateAll asynchronously locates all provided GeoLocs. And returns any
// errors that occur. An empty array on success.
func LocateAll(ctx context.Context, locater *Locater, unlocated []*models.GeoLoc) []error {
	// Save any errors that occur
	errs := []error{}

	// Channels to drive async jobs
	errsChan := make(chan error)
	locs := make(chan *models.GeoLoc)

	// Count of how many jobs are still running
	currentlyLocating := len(unlocated)

	// Kick off locator jobs
	for _, loc := range unlocated {
		// Locate
		locater.LocateAsync(ctx, errsChan, locs, loc)
	}

	// Wait for jobs to finish
	for currentlyLocating > 0 {
		select {
		case err := <-errsChan:
			// If error
			currentlyLocating -= 1
			errs = append(errs, fmt.Errorf("error locating model:"+
				" %s", err.Error()))
		case <-locs:
			// If success
			currentlyLocating -= 1
		}
	}

	// Return errs
	return errs
}
