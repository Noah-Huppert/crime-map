package models

const (
	// StatusOk indicates that a GAPI request was successful
	StatusOk string = "OK"

	// StatusZero indicates that the GAPI request didn't return any results
	StatusZero string = "ZERO_RESULTS"

	// StatusLimit indicates that we have gone over the GAPI query limit
	StatusLimit string = "OVER_QUERY_LIMIT"

	// StatusDenied indicates that the GAPI request was denied
	StatusDenied string = "REQUEST_DENIED"

	// StatusInvalid indicates that the GAPI request was invalid
	StatusInvalid string = "INVALID_REQUEST"

	// StatusErr indicates that an unknown error occurred during the GAPI
	// request
	StatusErr string = "UNKNOWN_ERROR"
)

const (
	// AccuracyPerfect indicates that the location provided by the GAPI
	// is exact
	AccuracyPerfect string = "ROOFTOP"

	// AccuracyBetween indicates that the location provided by the GAPI
	// is between two addresses
	AccuracyBetween string = "RANGE_INTERPOLATED"

	// AccuracyCenter indicates that the location is in the middle of an
	// region. Such as a block
	AccuracyCenter string = "GEOMETRIC_CENTER"

	// AccuracyApprox indicates that the location is not exact
	AccuracyApprox string = "APPROXIMATE"
)

// GeoLoc holds information about the geographical location of a crime "location"
// field.
//
// GeoLoc are resolved by using the Google Geocoding API to transform location
// strings into lat long coords. GeoLoc also holds some additional accuracy
// information, as not all locations can be resolved exactly.
type GeoLoc struct {
	Model

	// Located indicates if the raw location has been geocoded using the
	// GAPI
	Located bool

	// Lat is the latitude of the location
	Lat float64

	// Long is the longitude of the location
	Long float64

	// PostalAddr holds the formatted postal address of the location
	PostalAddr string

	// Accuracy indicates how close to the provided location the lat long
	// are
	Accuracy string

	// Partial indicates if the match is only a partial
	Partial bool

	// NeViewLat holds the recommended latitude which the northeast corner
	// of the map viewport should be located at
	NeViewLat float64

	// NeViewLong holds the recommended longitude which the northeast corner
	// of the map viewport should be located at
	NeViewLong float64

	// SwViewLat holds the recommended latitude which the southwest corner
	// of the map viewport should be located at
	SwViewLat float64

	// SwViewLong holds the recommended longitude which the southwest corner
	// of the map viewport should be located at
	SwViewLong float64

	// BoundsProvided indicates whether any location bounds were provided
	BoundsProvided bool

	// NeBoundsLat holds the location bounds northeast latitude
	NeBoundsLat float64

	// NeBoundsLong holds the location bounds northeast longitude
	NeBoundsLong float64

	// SwBoundsLat holds the location bounds southwest latitude
	SwBoundsLat float64

	// SwBoundsLong holds the location bounds southwest longitude
	SwBoundsLong float64

	// GAPIPlaceID holds the GAPI location ID, used to retrieve additional
	// information about a location using the GAPI
	GAPIPlaceID string

	// GAPIStatus holds the status of the GAPI Geocoding request which was
	// sent to locate the location
	GAPIStatus string

	// Raw holds the text present on the crime report which the GeoLoc
	// attempts to locate
	Raw string
}

// NewGeoLoc returns a new GeoLoc instance with the provided raw text
func NewGeoLoc(raw string) *GeoLoc {
	return &GeoLoc{
		Parsed: false,
		Raw:    raw,
	}
}
