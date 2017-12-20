package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pdf "github.com/unidoc/unidoc/pdf/model"
)

// TODO: Document code more

const file = "data/2017-10-12.pdf"

func main() {
	fmt.Printf("Input file: %s\n", file)
	err := inspectPdf(file)
	if err != nil {
		fmt.Printf("error inspecting file: %s\n", err.Error())
		os.Exit(1)
	}
}

// TODO: Make fn return a list of crimes in pdf
// TODO: Break up inspectPdf into smaller fns
func inspectPdf(inputPath string) error {
	// Open file
	pdfFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening pdf file: %s", err.Error())
	}

	defer pdfFile.Close()

	// Create pdf reader for file
	pdfReader, err := pdf.NewPdfReader(pdfFile)
	if err != nil {
		return fmt.Errorf("error creating pdf reader: %s", err.Error())
	}

	// Get number of pages
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return fmt.Errorf("error getting number of pages in pdf: %s",
			err.Error())
	}

	// Loop through pages
	rawFields := []string{}
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		// Get page
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return fmt.Errorf("error getting page #%d: %s",
				pageNum, err.Error())
		}

		// Get streams
		streams, err := page.GetAllContentStreams()
		if err != nil {
			return fmt.Errorf("error getting content streams: %s", err.Error())
		}

		inField := false
		field := ""
		for _, stream := range streams {
			s := string(stream)

			if s == "(" {
				inField = true
			} else if inField && s == ")" {
				inField = false
				rawFields = append(rawFields, field)
				field = ""
			} else if inField {
				field += s
			}

		}
	}

	// Loop through fields and group by crime
	crimesFields := [][]string{}
	wrkingFields := []string{}
	inPageHeader := false

	for _, field := range rawFields {
		// Start of a new crime
		if field == "Location :" {
			if len(wrkingFields) > 0 {
				crimesFields = append(crimesFields, wrkingFields)
			}

			wrkingFields = []string{}

			fmt.Println("")
		} else {
			// If midway through parsing a crime's fields

			// Check if page number
			if _, err := strconv.Atoi(field); err == nil {
				inPageHeader = true
				continue
			}

			// Check if header val
			split := strings.Split(field, " ")
			if field == "Student Right To Know Case Log Daily Report" {
				inPageHeader = false
				continue
			} else if len(split) > 0 && split[0] == "From" {
				inPageHeader = true
				continue
			} else if inPageHeader {
				continue
			}

			// Check if ignored
			if field == "Page No." ||
				field == " Report #:" ||
				field == "Date and Time Occurred From - Occurred To:" ||
				field == "Print Date and Time" ||
				field == "at" ||
				field == "Incident\\s\\" ||
				field == "Date Reported:" ||
				field == "Disposition:" ||
				field == "Synopsis:" {
				continue
			}

			// If total count of crimes
			if field == fmt.Sprintf(" %d", len(crimesFields)+1) {
				continue
			}

			// Add to working fields array
			wrkingFields = append(wrkingFields, field)
		}
	}

	// Add last wrkingFields result if not empty
	// But ignore last field, b/c this will be the total count of crimes
	if len(wrkingFields) > 0 {
		crimesFields = append(crimesFields, wrkingFields)
	}

	// Loop through crime fields and transform into Crime structs
	crimes := []Crime{}

	for i, fields := range crimesFields {
		// Check at least 5 fields provided
		numFields := len(fields)
		if numFields < 5 {
			return fmt.Errorf("error parsing crime: not enough "+
				"fields, found %d, needs 5, %s", len(fields),
				fields)
		}

		// Set fields
		crime := Crime{}

		// TODO: Parse into more meaningful data types
		crime.DateReported = fields[0]
		crime.Location = fields[1]
		crime.ReportID = fields[2]
		crime.Incidents = fields[3]
		crime.DateOccurred = fields[4]

		// If description provided
		if numFields > 5 {
			crime.Descriptions = fields[4 : numFields-2]
		}

		crime.Remediation = fields[numFields-1]

		fmt.Printf("\n%d\n====\n%s\n", i+1, crime)

		crimes = append(crimes, crime)
	}

	fmt.Printf("%d crimes\n", len(crimesFields))

	return nil
}

type Crime struct {
	DateReported string
	DateOccurred string
	ReportID     string
	Location     string
	Incidents    string
	Descriptions []string
	Remediation  string
}

func (c Crime) String() string {
	return fmt.Sprintf("Reported: %s\n"+
		"Occurred: %s\n"+
		"ID: %s\n"+
		"Location: %s\n"+
		"Incidents: %s\n"+
		"Description: %s\n"+
		"Remediation: %s",
		c.DateReported, c.DateOccurred, c.ReportID, c.Location,
		c.Incidents, strings.Join(c.Descriptions, ","), c.Remediation)
}
