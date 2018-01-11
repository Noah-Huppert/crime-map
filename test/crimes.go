package test

import (
	"testing"

	"github.com/Noah-Huppert/crime-map/models"
)

// CrimesSlicesEq determines if a slice of Crime structs are equal
func CrimesSlicesEq(t *testing.T, expected []models.Crime, actual []models.Crime) {
	// Check length
	if len(expected) != len(actual) {
		t.Fatalf("expected length not equal to actual length, "+
			"len(expected) = %d, len(actual) = %d", len(expected),
			len(actual))
	}

	// Loop through and check Crimes
	for i, expVal := range expected {
		actVal := actual[i]

		if !(expVal.Equal(actVal)) {
			t.Fatalf("%d index value of actual does not equal "+
				"expected, expected[i]: %s, actual[i]: %s",
				i, expVal, actVal)
		}
	}
}
