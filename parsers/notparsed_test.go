package parsers

import (
	"testing"
)

// expectedNotParsedErrStr is the expected string to be returned in
// TestNotParsedString when the ErrFieldNotParsed.Error() method is called.
var expectedNotParsedErrStr = "field not parsed, index: 5, field: \"A\""

// TestNotParsedString ensures the ErrFieldNotParsed.Error method returns the
// correct value for the constructed error.
func TestNotParsedString(t *testing.T) {
	e := NewErrFieldNotParsed(5, "A")

	// Test
	if e.Error() != expectedNotParsedErrStr {
		t.Fatalf("error did not matched expected value, \n"+
			"expected: \"%s\",\n"+
			"actual  : \"%s\"", expectedNotParsedErrStr, e.Error())
	}
}
