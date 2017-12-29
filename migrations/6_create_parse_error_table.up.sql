CREATE TABLE parse_errors (
	id SERIAL PRIMARY KEY,

	crime_id INTEGER REFERENCES crimes,

	field TEXT NOT NULL,
	original TEXT NOT NULL,
	corrected TEXT NOT NULL,

	err_type ERR_TYPE_T NOT NULL
)
