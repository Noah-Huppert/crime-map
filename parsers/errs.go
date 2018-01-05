package parsers

import (
	"errors"
)

// ErrReportParsed is the error returned when a Parser has already parsed a
// report
var ErrReportParsed error = errors.New("report already parsed")

// ErrReportNotParsed is the error returned when a Parser has not been parsed
// yet
var ErrReportNotParsed error = errors.New("report has not been parsed")
