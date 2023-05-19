package main

import "testing"

func TestInStoplist(t *testing.T) {
	tests := []struct {
		word     string
		expected bool
	}{
		{"the", true},
		{"don't", true},
		{"dont", false},
		{"hippopotamus", false},
	}

	for _, tc := range tests {
		got := InStoplist(tc.word)

		if got != tc.expected {
			t.Errorf("got %t, wanted %t", got, tc.expected)
		}
	}
}
