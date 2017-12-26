package models

// GeoBound indicates a square area on a map
type GeoBound struct {
	// ID is the unique identifier
	ID uint

	// NeLat holds the recommended latitude which the northeast corner
	// of the map viewport should be located at
	NeLat float32

	// NeLong holds the recommended longitude which the northeast corner
	// of the map viewport should be located at
	NeLong float32

	// SwLat holds the recommended latitude which the southwest corner
	// of the map viewport should be located at
	SwLat float32

	// SwLong holds the recommended longitude which the southwest corner
	// of the map viewport should be located at
	SwLong float32
}
