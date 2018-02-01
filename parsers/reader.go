package parsers

import (
	"fmt"
	"time"

	"github.com/Noah-Huppert/crime-map/errs"
	"github.com/Noah-Huppert/crime-map/geo"
	"github.com/Noah-Huppert/crime-map/models"
	"github.com/Noah-Huppert/crime-map/pdf"
)

// Reader takes in a Pdf file, and extracts crimes from it. Using the
// appropriate parser, based on the header.
type Reader struct {
	// geoCache is used to cache GeoLoc queries
	geoCache *geo.GeoCache
}

// NewReader creates a new Reader struct with the given file path.
func NewReader(geoCache *geo.GeoCache) *Reader {
	return &Reader{
		geoCache: geoCache,
	}
}

/*
// IsParsed indicates if the specified pdf file has been parsed for crimes yet
func (r Reader) IsParsed() bool {
	return r.parsed
}

// Crimes returns the crimes recorded in the specified Clery report. Along with
// a boolean indicating if the report has been parsed yet
func (r Reader) Crimes() ([]models.Crime, bool) {
	return r.crimes, r.IsParsed()
}
*/

// Parse interprets a crime report file and returns the contained crimes.
// Additionally an error will be returned, nil on success.
func (r *Reader) Parse(path string) ([]*models.Crime, error) {
	// Pre parse report information
	file, report, err := r.preParseReport(path)
	if err != nil {
		return nil, fmt.Errorf("error gathering pre parse report "+
			"information: %s", err.Error())
	}

	// Parse report
	crimes, err := r.parseReport(file, report)
	if err != nil {
		return nil, fmt.Errorf("error parsing report: %s", err.Error())
	}

	// Save information about parsing process itself in Report model
	crimesCount := uint(len(crimes))
	err = r.updateReportPost(crimesCount, report)
	if err != nil {
		return nil, fmt.Errorf("error updating report model after"+
			" parsing: %s", err.Error())
	}

	// All done
	return crimes, nil
}

// preParseReport determines some preliminary information about the pdf report
// being parsed. Such as the date range covered and the number of pages.
//
// If a Report model is found in the database it will be queried. If not one
// will be inserted. The Report model will be returned, with a populated
// Report.ID field.
//
// Additionally an error will be returned if one occurs, nil on success.
func (r Reader) preParseReport(path string) (*pdf.Pdf, *models.Report, error) {
	// Open pdf
	file := pdf.NewPdf(path)

	// Parse PDF
	if _, err := file.Parse(); err != nil {
		return nil, nil, fmt.Errorf("error parsing pdf: %s", err.Error())
	}

	// Time parsed
	report := &models.Report{}
	now := time.Now()
	report.ParsedOn = &now

	// Determine number of pdf report pages
	pages, parsed := file.Pages()
	if !parsed {
		return nil, nil, errs.ErrNotParsed
	}
	report.Pages = pages

	// Get report fields to parse addition metadata
	fields, parsed := file.Fields()
	if !parsed {
		return nil, nil, errs.ErrNotParsed
	}

	// Determine university
	univ, err := determineUniversity(fields)
	if err != nil {
		return nil, nil, fmt.Errorf("error determining university "+
			"from report fields: %s", err.Error())
	}
	report.University = univ

	// Determine date range
	dateRunner := NewOnceRunner(DateRangeParser{})
	err = dateRunner.Parse(report, nil, fields)

	if err != nil {
		return nil, nil, fmt.Errorf("error parsing report date range: %s",
			err.Error())
	}

	// Save report
	if err = report.InsertIfNew(); err != nil {
		return nil, nil, fmt.Errorf("error saving Report model: %s",
			err.Error())
	}

	// Success
	return file, report, nil
}

func (r Reader) parseReport(file *pdf.Pdf, report *models.Report) ([]*models.Crime, error) {
	// Use parser based on university
	var runner *ParserRunner

	if report.University == models.UniversityDrexel {
		runner = NewDrexelRunner()
	} else {
		return nil, fmt.Errorf("no parser parser for university:"+
			" %s", report.University)
	}

	// Get pdf text fields
	fields, err := file.Parse()
	if err != nil {
		return nil, fmt.Errorf("error getting pdf fields: %s\n", err.Error())
	}

	// Check if report has already been parsed
	if report.ParseSuccess {
		// Skip parsing if already parsed
		return nil, errs.ErrParsed
	}

	// Parse crimes from fields
	crimes, err := runner.Parse(report, fields)
	if err != nil {
		return nil, fmt.Errorf("error parsing report: %s",
			err.Error())
	}

	// Success
	return crimes, nil
}

// updateReportPost sets the ParseSuccess and CrimesCount properties of the
// Report model associated with the parsing job.
func (r Reader) updateReportPost(count uint, report *models.Report) error {
	// Get number of crimes parsed
	report.CrimesCount = count

	// Indicate report parsed successfully
	report.ParseSuccess = true

	// Save updates
	err := report.UpdatePostParseFields()
	if err != nil {
		return fmt.Errorf("error saving post parse updates to report "+
			"model: %s", err.Error())
	}

	// Success
	return nil
}
