CREATE TABLE geo_locs (
	id SERIAL PRIMARY KEY,

	located BOOLEAN NOT NULL DEFAULT FALSE,
	gapi_success BOOLEAN NOT NULL DEFAULT FALSE,

	lat FLOAT(32),
	long FLOAT(32),

	postal_addr TEXT,

	accuracy GEO_ACCURACY_T,

	partial BOOLEAN DEFAULT FALSE,
	bounds_provided BOOLEAN DEFAULT FALSE,

	bounds_id INTEGER REFERENCES geo_bounds,

	gapi_place_id TEXT,

	raw TEXT NOT NULL UNIQUE
)
