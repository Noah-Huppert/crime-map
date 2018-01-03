package gapi

import (
	"fmt"
	"googlemaps.github.io/maps"

	"github.com/Noah-Huppert/crime-map/config"
)

// client holds the Google API client if retrieved, nil if not
var client *maps.Client

// NewClient creates a new Google API client with the credentials from the
// configuration file. An error is returned if one occurs, or nil on success.
func NewClient() (*maps.Client, error) {
	// Check if we have client
	if client != nil {
		return client, nil
	}

	// Get config
	config, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("error retrieving configuration: %s",
			err.Error())
	}

	// Make client
	client, err = maps.NewClient(maps.WithAPIKey(config.GAPI.APIKey))
	if err != nil {
		return nil, fmt.Errorf("error creating GAPI client: %s",
			err.Error())
	}

	// Success
	return client, nil
}
