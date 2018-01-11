package test

import "testing"

// StrSlicesEq checks that the two provided slices have equal elements in
// the same order.
func StrSlicesEq(t *testing.T, expected []string, actual []string) {
	// Check length
	if len(expected) != len(actual) {
		t.Fatalf("slices have different lengths, len expected: %d, "+
			"len actual: %d", len(expected), len(actual))
	}

	// Check equal and order
	for i, valExp := range expected {
		valAct := actual[i]

		if valExp != valAct {
			t.Fatalf("slice values at index %d are not equal, "+
				"expected: %#v, b: %#v", i, valExp, valAct)
		}
	}
}
