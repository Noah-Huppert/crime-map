package parsers

import (
	"errors"
	"fmt"
	"github.com/Noah-Huppert/crime-map/geo"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/pdf"
	"time"
)

// Report structs hold information about Clery act crime log pdf files
type Report struct {
	// pdf holds the Pdf object used to parse the report file
	pdf *pdf.Pdf

	// University holds the name of the school which the crime reports
	// were collected from
	University string

	// StartDate is the start of the date range which the report records
	// crimes for
	StartDate time.Time

	// EndDate is the end of the date range which the report records
	// crimes for
	EndDate time.Time

	// parsed indicates whether the pdf's text fields have been converted
	// into Crime structs
	parsed bool

	// crimes holds all the crimes found in the clery report
	crimes []models.Crime

	// geoCache is used to cache GeoLoc queries
	geoCache *geo.GeoCache
}

// NewReport creates a new report struct with the given file path. Additionally
// an error is returned if one occurs, or nil on success.
func NewReport(path string, geoCache *geo.GeoCache) *Report {
	return &Report{
		pdf:      pdf.NewPdf(path),
		parsed:   false,
		crimes:   []models.Crime{},
		geoCache: geoCache,
	}
}

// IsParsed indicates if the specified pdf file has been parsed for crime yet
func (r Report) IsParsed() bool {
	return r.parsed
}

// Crimes returns the crimes recorded in the specified Clery report. Along with
// a boolean indicating if the report has been parsed yet
func (r Report) Crimes() ([]models.Crime, bool) {
	return r.crimes, r.IsParsed()
}

// Parse interprets a crime report file and returns the contained crimes.
// Additionally an error will be returned, nil on success.
func (r Report) Parse() ([]models.Crime, error) {
	// Check if parsed
	if r.IsParsed() {
		return r.crimes, errors.New("report has already been parsed")
	}

	// Get pdf text fields
	fields, err := r.pdf.Parse()
	if err != nil {
		fmt.Printf("error getting pdf fields: %s\n", err.Error())
	}

	// Parse into crimes
	p := NewDrexelParser(r.geoCache)

	crimes, err := p.Parse(fields)
	if err != nil {
		fmt.Printf("error parsing report: %s", err.Error())
	}

	// All done
	return crimes, nil
}
