package errs

import (
	"errors"
)

// ErrParsed indicates that the provided data was already parsed
var ErrParsed error = errors.New("data already parsed")

// ErrNotParsed indicates that the provided data was not parsed by any Parsers
var ErrNotParsed error = errors.New("no data parsed")

// ErrCrimeParsed indicates that the provided crime has been completely
// parsed during the invocation of the parser.
//
// If this error is received: Append the crime you are currently providing
// to a list of parsed crimes. Then create a new empty crime model. And provide
// this on the next invocation of the parser.
var ErrCrimeParsed error = errors.New("crime has been successfully parsed")
