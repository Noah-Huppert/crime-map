package parsers

import (
	"testing"

	"github.com/Noah-Huppert/crime-map/errs"
	"github.com/Noah-Huppert/crime-map/models"
)

// TestFindsDrexel ensures the determineUniversity function identifies Drexel
// University from a list of fields.
func TestFindsDrexel(t *testing.T) {
	univ, err := determineUniversity([]string{"not the field", "or this",
		"or me", string(models.UniversityDrexel), "some at the end",
		"single", "words"})

	// Test
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if univ != models.UniversityDrexel {
		t.Fatalf("university did not match expected: %s, actual: %s",
			models.UniversityDrexel, univ)
	}
}

// TestUnknownUnivErrs ensures the determineUniversity function returns an
// error if no known university was found
func TestUnknownUnivErrs(t *testing.T) {
	_, err := determineUniversity([]string{"not", "a", "univ"})

	// Test
	if (err == nil) || (err != errs.ErrUnknownUniv) {
		t.Fatalf("error did not match expected: %s, actual: %s",
			errs.ErrUnknownUniv, err)
	}
}
