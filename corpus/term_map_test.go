package corpus

import (
	"testing"
)

func TestTermID(t *testing.T) {
	// Clear out the existing cache mapping terms to IDs
	termIDs = make(map[string]int)

	// Check that calling termID maps the next available int ID to the term and
	// fetches the ID associated with already-cached terms
	tests := []struct {
		term string
		id   int
	}{
		{"foo", 0},
		{"bar", 1},
		{"foo", 0},
		{"baz", 2},
	}

	for _, tc := range tests {
		id := termID(tc.term)

		if id != tc.id {
			t.Errorf("got %d, wanted %d", tc.id, id)
		}
	}
}
