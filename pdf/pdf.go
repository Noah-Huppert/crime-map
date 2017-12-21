package pdf

import (
	"errors"
	"fmt"
	pdfcontent "github.com/unidoc/unidoc/pdf/contentstream"
	pdfcore "github.com/unidoc/unidoc/pdf/core"
	pdf "github.com/unidoc/unidoc/pdf/model"
	"os"
)

// Pdf holds data about a pdf file
type Pdf struct {
	// path is the location of the pdf file
	path string

	// parsed indicates whether or not fields have been extracted from the
	// pdf file
	parsed bool

	// fields holds all the text fields present in the pdf, in the order
	// they occurred
	fields []string
}

// NewPdf creates a new Pdf struct with the given path
func NewPdf(path string) *Pdf {
	return &Pdf{
		path:   path,
		parsed: false,
		fields: []string{},
	}
}

// IsParsed indicates if the specified pdf file has been processed yet
func (p Pdf) IsParsed() bool {
	return p.parsed
}

// Fields returns the fields the Pdf contains. Along with a boolean, which
// indicates if the pdf file has been parsed yet.
func (p Pdf) Feilds() ([]string, bool) {
	return p.fields, p.IsParsed()
}

// Parse opens the pdf file and extracts all text fields present. These fields
// are returned. Along with an error if one occurs, or nil on success.
func (p Pdf) Parse() ([]string, error) {
	// If already parsed, error
	if p.IsParsed() {
		return p.fields, errors.New("pdf file already parsed")
	}

	// Open file
	file, err := os.Open(p.path)
	if err != nil {
		return p.fields, fmt.Errorf("error opening pdf file: %s", err.Error())
	}

	defer file.Close()

	// Create pdf reader for file
	pdfReader, err := pdf.NewPdfReader(file)
	if err != nil {
		return p.fields, fmt.Errorf("error creating pdf reader: %s", err.Error())
	}

	// Determine if pdf is encrypted
	isEncrypted, err := pdfReader.IsEncrypted()
	if err != nil {
		return p.fields, fmt.Errorf("error determining pdf file "+
			"encryption status: %s", err.Error())
	}

	// If encrypted, try decrypting with empty password
	if isEncrypted {
		// Attempt
		auth, err := pdfReader.Decrypt([]byte(""))
		if err != nil {
			return p.fields, fmt.Errorf("error decrypting pdf "+
				"file: %s", err.Error())
		}

		// Verify successful decryption
		if !auth {
			return p.fields, errors.New("unable to decrypt pdf " +
				"file with empty password")
		}
	}

	// Get number of pages
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return p.fields, fmt.Errorf("error getting number of pages in pdf: %s",
			err.Error())
	}

	// Loop through pages
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		// Get page
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return p.fields, fmt.Errorf("error getting pdf page #%d: %s",
				pageNum, err.Error())
		}

		// Get page contents
		streams, err := page.GetAllContentStreams()
		if err != nil {
			return p.fields, fmt.Errorf("error getting pdf content streams: %s",
				err.Error())
		}

		// Parse contents
		parser := pdfcontent.NewContentStreamParser(streams)
		ops, err := parser.Parse()
		if err != nil {
			return p.fields, fmt.Errorf("error parsing content "+
				"stream: %s", err.Error())
		}

		// Loop through text
		for _, op := range *ops {
			// Check text field
			if op.Operand == "Tj" && len(op.Params) == 1 {
				val, ok := op.Params[0].(*pdfcore.PdfObjectString)
				if !ok {
					return p.fields, errors.New("error " +
						"casting pdf text field to string")
				}

				p.fields = append(p.fields, string(*val))
			}
		}
	}

	// Done
	p.parsed = true
	return p.fields, nil
}
