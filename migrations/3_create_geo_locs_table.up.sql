CREATE TABLE geo_locs (
	id SERIAL PRIMARY KEY,

	located BOOLEAN NOT NULL,

	lat FLOAT(32),
	long FLOAT(32),

	postal_addr TEXT,

	accuracy GEO_ACCURACY_T,

	partial BOOLEAN,
	bounds_provided BOOLEAN,

	bounds_id INTEGER REFERENCES geo_bounds,

	gapi_place_id TEXT,

	raw TEXT NOT NULL
)
