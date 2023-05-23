package main

import "testing"

func TestDefaultStoplist(t *testing.T) {
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
		got := DefaultStoplist.Include(tc.word)

		if got != tc.expected {
			t.Errorf("got %t, wanted %t", got, tc.expected)
		}
	}
}

func TestCustomStoplist(t *testing.T) {
	stoplist := NewStoplist([]string{
		"foo",
		"bar",
	})

	tests := []struct {
		word     string
		expected bool
	}{
		{"foo", true},
		{"bar", true},
		{"baz", false},
	}

	for _, tc := range tests {
		got := stoplist.Include(tc.word)

		if got != tc.expected {
			t.Errorf("got %t, wanted %t", got, tc.expected)
		}
	}
}
