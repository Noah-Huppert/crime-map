CREATE TABLE geo_bounds (
	id SERIAL PRIMARY KEY,

	ne_lat FLOAT(32) NOT NULL,
	ne_long FLOAT(32) NOT NULL,

	sw_lat FLOAT(32) NOT NULL,
	sw_long FLOAT(32) NOT NULL
)
