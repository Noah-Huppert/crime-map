CREATE TABLE geo_locs (
	id SERIAL PRIMARY KEY,

	located BOOLEAN NOT NULL,

	lat FLOAT(32) NOT NULL,
	long FLOAT(32) NOT NULL,

	postal_addr TEXT NOT NULL,

	accuracy GEO_ACCURACY_T NOT NULL,


)
