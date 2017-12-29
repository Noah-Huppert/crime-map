CREATE TABLE crimes (
	id SERIAL PRIMARY KEY,

	date_reported TIMESTAMP WITH TIME ZONE NOT NULL,
	date_occurred TSTZRANGE NOT NULL,

	report_super_id INTEGER NOT NULL,
	report_sub_id INTEGER NOT NULL,

	location TEXT NOT NULL,
	geo_loc_id INTEGER REFERENCES geo_locs,

	incidents TEXT[] NOT NULL,
	descriptions TEXT[] NOT NULL,

	remediation TEXT
)