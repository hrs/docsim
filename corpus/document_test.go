package corpus

import (
	"reflect"
	"testing"
)

func TestNormalizeTfIdf(t *testing.T) {
	tm := termMap{
		"foo": 2.0,
		"bar": 3.0,
		"baz": 4.0,
	}

	tests := []struct {
		doc    Document
		idf    termMap
		expDoc Document
	}{
		{
			Document{
				termFreq: termMap{},
			},
			tm,
			Document{
				tfIdf: termMap{},
				norm:  0.0,
			},
		},
		{
			Document{
				termFreq: termMap{
					"foo": 3.0,
					"bar": 4.0,
					"baz": 5.0,
				},
			},
			tm,
			Document{
				tfIdf: termMap{
					"foo": 6.0,
					"bar": 12.0,
					"baz": 20.0,
				},
				norm: 24.0832,
			},
		},
	}

	for _, tc := range tests {
		tc.doc.normalizeTfIdf(tc.idf)

		if !reflect.DeepEqual(tc.doc.tfIdf, tc.expDoc.tfIdf) {
			t.Errorf("got %v, wanted %v", tc.doc.tfIdf, tc.expDoc.tfIdf)
		}

		if !approxEq(tc.doc.norm, tc.expDoc.norm) {
			t.Errorf("got %.4f, wanted %.4f", tc.doc.norm, tc.expDoc.norm)
		}

		if tc.doc.termFreq != nil {
			t.Errorf("got %v, wanted nil", tc.doc.termFreq)
		}
	}
}
func TestCalcNorm(t *testing.T) {
	tests := []struct {
		doc      Document
		expected float64
	}{
		{
			Document{
				tfIdf: termMap{},
			},
			0.0,
		},
		{
			Document{
				tfIdf: termMap{
					"foo": 2.0,
					"bar": 3.0,
					"baz": 4.0,
				},
			},
			5.3852,
		},
	}

	for _, tc := range tests {
		got := tc.doc.calcNorm()

		if !approxEq(got, tc.expected) {
			t.Errorf("got %.4f, wanted %.4f", got, tc.expected)
		}
	}
}
