package parsers

import (
	"fmt"
	"time"

	"github.com/Noah-Huppert/crime-map/geo"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/pdf"
)

// Reader takes in a Pdf file, and extracts crimes from it. Using the
// appropriate parser, based on the header.
type Reader struct {
	// pdf holds the Pdf object used to parse the report file
	pdf *pdf.Pdf

	// parsed indicates whether the pdf's text fields have been converted
	// into Crime structs
	parsed bool

	// crimes holds all the crimes found in the clery report
	crimes []models.Crime

	// geoCache is used to cache GeoLoc queries
	geoCache *geo.GeoCache
}

// NewReader creates a new Reader struct with the given file path.
func NewReader(path string, geoCache *geo.GeoCache) *Reader {
	return &Reader{
		pdf:      pdf.NewPdf(path),
		parsed:   false,
		crimes:   []models.Crime{},
		geoCache: geoCache,
	}
}

// IsParsed indicates if the specified pdf file has been parsed for crimes yet
func (r Reader) IsParsed() bool {
	return r.parsed
}

// Crimes returns the crimes recorded in the specified Clery report. Along with
// a boolean indicating if the report has been parsed yet
func (r Reader) Crimes() ([]models.Crime, bool) {
	return r.crimes, r.IsParsed()
}

// Parse interprets a crime report file and returns the contained crimes.
// Additionally an error will be returned, nil on success.
func (r Reader) Parse() ([]models.Crime, error) {
	// Check if parsed
	if r.IsParsed() {
		return r.crimes, ErrReportParsed
	}

	// Get pdf text fields
	fields, err := r.pdf.Parse()
	if err != nil {
		return r.crimes, fmt.Errorf("error getting pdf fields: %s\n", err.Error())
	}

	// Figure out which university published report
	univ, err := determineUniversity(fields)
	if err != nil {
		return r.crimes, fmt.Errorf("error determining university "+
			"from report fields: %s", err.Error())
	}

	// Use parser based on university
	var parser Parser

	if univ == models.UniversityDrexel {
		parser = NewDrexelParser(r.geoCache, fields)
	} else {
		return r.crimes, fmt.Errorf("error determining parser based on"+
			" university, no parser, university: %s", univ)
	}

	// Save Report model based on info in pdf
	report, err := r.saveReport(parser, univ)
	if err != nil {
		return r.crimes, fmt.Errorf("error saving report model: %s",
			err.Error())
	}

	// Parse crimes from fields
	crimes, err := parser.Parse(report.ID)
	if err != nil {
		return r.crimes, fmt.Errorf("error parsing report: %s",
			err.Error())
	}
	r.crimes = crimes

	// All done
	return r.crimes, nil
}

// saveReport retrieves information about the report being parsed, and
// retrieves / inserts a report with the information. An error is returned if
// one occurs, nil on success.
func (r Reader) saveReport(parser Parser, univ models.UniversityType) (*models.Report, error) {
	// Get date range report covers
	startRange, endRange, err := parser.Range()
	if err != nil {
		return nil, fmt.Errorf("error retrieving report range: %s",
			err.Error())
	}

	// Get number of pages in report
	pages, parsed := r.pdf.Pages()
	if !parsed {
		return nil, fmt.Errorf("error retrieving number of report"+
			" pages: %s", ErrReportNotParsed)
	}

	// Make report
	now := time.Now()
	report := models.NewReport(univ, &now, startRange, endRange,
		pages)

	// Save report
	if err = report.InsertIfNew(); err != nil {
		return nil, fmt.Errorf("error saving Report model: %s",
			err.Error())
	}

	// Success
	return report, nil
}
