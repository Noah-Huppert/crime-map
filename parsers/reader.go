package parsers

import (
	"errors"
	"fmt"
	"strings"

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
		return r.crimes, errors.New("report has already been parsed")
	}

	// Get pdf text fields
	fields, err := r.pdf.Parse()
	if err != nil {
		return r.crimes, fmt.Errorf("error getting pdf fields: %s\n", err.Error())
	}

	// Figure out which university published report
	univ, err := r.determineUniversity(fields)
	if err != nil {
		return r.crimes, fmt.Errorf("error determining university "+
			"from report fields: %s", err.Error())
	}

	// Parse into crimes
	drexelParser := NewDrexelParser(r.geoCache)

	crimes, err := drexelParser.Parse(fields)
	if err != nil {
		return r.crimes, fmt.Errorf("error parsing report: %s",
			err.Error())
	}

	// Get date range report covers
	startRange, endRange, err := drexelParser.Range()
	if err != nil {
		return r.crimes, fmt.Errorf("error retrieving report range: %s",
			err.Error())
	}

	// Get number of pages in report
	pages, parsed := r.pdf.Pages()
	if !parsed {
		return r.crimes, fmt.Errorf("error retrieving number of report"+
			" pages: %s", ErrReportNotParsed)
	}

	// Make report
	report := models.NewReport(univ, *startRange, *endRange, pages)

	// Save report
	if err = report.InsertIfNew(); err != nil {
		return crimes, fmt.Errorf("error saving Report model: %s",
			err.Error())
	}

	// Set models.Crime.ReportID fk in crimes
	for i, _ := range crimes {
		crimes[i].ReportID = report.ID
	}

	// All done
	return crimes, nil
}

// determineUniversity figures out which University a crime report was
// published from. By reading in the text fields present in a report. And
// searching for the first occurrence of a university name.
//
// A models.UniversityType is returned along with an error. Which will be nil
// on success.
func (r Reader) determineUniversity(fields []string) (models.UniversityType, error) {
	// Attempt to find univ name in fields
	for _, field := range fields {
		// Check
		if strings.Contains(field, string(models.UniversityDrexel)) {
			// Success
			return models.UniversityDrexel, nil
		}
	}

	// If none found
	return models.UniversityErr, errors.New("error determining university," +
		" no field with university name found")
}
