package config

import (
	"googlemaps.github.io/maps"
)

// GeoConfig holds configuration related to locating crimes from reports
type GeoConfig struct {
	// BoundsNeLat holds the northeast bounds latitude of the area to look
	// for crime locations in
	BoundsNeLat float64

	// BoundsNeLong holds the northeast bounds longitude of the area to
	// look for crime locations in
	BoundsNeLong float64

	// BoundsSwLat holds the southwest bounds latitude of the area to look
	// for crime locations in
	BoundsSwLat float64

	// BoundsSwLong holds the southwest bounds longitude of the area to
	// look for crime locations in
	BoundsSwLong float64

	// AddrPostfix is the string appended to the end of crime address
	// before attempting to locate it on a map
	AddrPostfix string
}

// MakeMapsBounds constructs a Google Maps map.LatLngBounds struct from the
// GeoConfig.Bounds* fields. This can be used to indicate the bounds to the
// Google Maps SDK.
func (c GeoConfig) MakeMapsBounds() *maps.LatLngBounds {
	return &maps.LatLngBounds{
		NorthEast: maps.LatLng{
			Lat: c.BoundsNeLat,
			Lng: c.BoundsNeLong,
		},
		SouthWest: maps.LatLng{
			Lat: c.BoundsSwLat,
			Lng: c.BoundsSwLong,
		},
	}
}
