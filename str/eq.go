package str

import "testing"

// SlicesEq checks that the two provided slices have equal elements in
// the same order.
//
// Equality errors will be printed with the provided testing.T object, if not
// nil.
func SlicesEq(t *testing.T, expected []string, actual []string) bool {
	// Check length
	if len(expected) != len(actual) {
		if t != nil {
			t.Fatalf("slices have different lengths, len expected: %d, "+
				"len actual: %d", len(expected), len(actual))
		}

		return false
	}

	// Check equal and order
	for i, valExp := range expected {
		valAct := actual[i]

		if valExp != valAct {
			if t != nil {
				t.Fatalf("slice values at index %d are not equal, "+
					"expected: %#v, b: %#v", i, valExp, valAct)
			}

			return false
		}
	}

	return true
}
