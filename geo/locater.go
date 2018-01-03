package geo

import (
	"context"
	"database/sql"
	"fmt"
	"googlemaps.github.io/maps"
	"strings"

	"github.com/Noah-Huppert/crime-map/config"
	"github.com/Noah-Huppert/crime-map/gapi"
	"github.com/Noah-Huppert/crime-map/models"
)

// region holds the ccTLD two-character value for the area where the GAPI
// should look for locations
const region string = "us"

// unknownLocRaw is the raw string value of the GeoLoc which indicates that a
// crime's location is unknown
const unknownLocRaw string = "UNKNOWN LOCATION - Non-reportable Location"

// Locater uses the Google Maps API to determine exactly where new GeoLoc
// models are in the world
type Locater struct{}

// NewLocater creates a new Locater instance
func NewLocater() *Locater {
	return &Locater{}
}

// Locate determines where a GeoLoc model resides on the map. Determining
// bounds and lat long. An error is returned if one occurs, or nil on success.
//
// A context must be provided to manage the GAPI request's running.
//
// If the GeoLoc provided is indicates the location is unknown, the method
// returns immediately.
func (l Locater) Locate(ctx context.Context, loc *models.GeoLoc) error {
	// Check if located
	if loc.Located {
		return fmt.Errorf("geoloc model already located")
	}

	// Check if unknown
	if loc.Raw == unknownLocRaw {
		// Just exit
		return nil
	}

	// Get api client
	client, err := gapi.NewClient()
	if err != nil {
		return fmt.Errorf("error retrieving GAPI client: %s",
			err.Error())
	}

	// Get configuration
	c, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("error retrieving configuration: %s",
			err.Error())
	}

	// Trim raw location string
	// Usually in form:
	// 	<actual addr> - <addr annotation>
	// So get rid of second part so geocoding API works better
	locStr := strings.Split(loc.Raw, " - ")[0]

	// Then get rid of any parenthesis as well
	locStr = strings.Split(locStr, " (")[0]

	// Add a postfix to the address to zero in on the area
	locStr += c.Geo.AddrPostfix

	// Construct Geocode request
	req := maps.GeocodingRequest{
		Address: locStr,
		Region:  region,
		Bounds:  c.Geo.MakeMapsBounds(),
	}

	// Make Geocode request
	res, err := client.Geocode(ctx, &req)
	if err != nil {
		// Indicate geocoding failed
		loc.GAPISuccess = false

		return fmt.Errorf("error geocoding location: %s", err.Error())
	}

	// Extract first/best result
	if len(res) == 0 {
		return fmt.Errorf("no geocoding results returned")
	}
	best := res[0]

	// Indicate geocoding request succeeded
	loc.GAPISuccess = true

	// Save results
	// Lat long
	loc.Lat = best.Geometry.Location.Lat
	loc.Long = best.Geometry.Location.Lng

	// Address
	loc.PostalAddr = best.FormattedAddress

	// Accuracy
	loc.Accuracy, err = models.NewGeoLocAccuracy(best.Geometry.LocationType)
	if err != nil {
		return fmt.Errorf("error parsing accuracy value: %s",
			err.Error())
	}

	// Bounds
	bounds := models.GeoBoundFromMapsBound(best.Geometry.Bounds)

	loc.BoundsProvided = (bounds.NeLat != 0) && (bounds.NeLong != 0) &&
		(bounds.SwLat != 0) && (bounds.SwLong != 0)

	// Insert bounds if provided
	if loc.BoundsProvided {
		if err = bounds.InsertIfNew(); err != nil {
			return fmt.Errorf("error querying/inserting location "+
				"bounds: %s", err.Error())
		}

		loc.BoundsID = sql.NullInt64{
			Int64: int64(bounds.ID),
			Valid: true,
		}
	}

	// Viewport bounds
	viewBounds := models.GeoBoundFromMapsBound(best.Geometry.Viewport)
	if err = viewBounds.InsertIfNew(); err != nil {
		return fmt.Errorf("error querying/inserting viewport bounds: %s",
			err.Error())
	}
	loc.ViewportBoundsID = viewBounds.ID

	// GAPI place ID
	loc.GAPIPlaceID = best.PlaceID

	// Indicate GeoLoc has been located
	loc.Located = true

	return nil
}

// LocateAsync wraps the Locate method with asynchronous logic. Passing errors
// through a provided channel, and recording finsihed work via a locs channel.
func (l Locater) LocateAsync(ctx context.Context, errs chan error,
	locs chan *models.GeoLoc, loc *models.GeoLoc) {

	// Start async
	go func() {
		// Locate
		err := l.Locate(ctx, loc)

		// If error
		if err != nil {
			errs <- fmt.Errorf("error running async Locate, loc: "+
				"%s, err: %s", loc, err.Error())
		} else {
			// If success
			locs <- loc
		}
	}()
}
