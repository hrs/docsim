package main

import (
	"reflect"
	"testing"
)

func TestNormalizeTfIdf(t *testing.T) {
	tm := TermMap{
		"foo": 2.0,
		"bar": 3.0,
		"baz": 4.0,
	}

	tests := []struct {
		doc    Document
		idf    TermMap
		expDoc Document
	}{
		{
			Document{
				TermFreq: TermMap{},
			},
			tm,
			Document{
				TfIdf: TermMap{},
				Norm:  0.0,
			},
		},
		{
			Document{
				TermFreq: TermMap{
					"foo": 3.0,
					"bar": 4.0,
					"baz": 5.0,
				},
			},
			tm,
			Document{
				TfIdf: TermMap{
					"foo": 6.0,
					"bar": 12.0,
					"baz": 20.0,
				},
				Norm: 24.0832,
			},
		},
	}

	for _, tc := range tests {
		tc.doc.NormalizeTfIdf(tc.idf)

		if !reflect.DeepEqual(tc.doc.TfIdf, tc.expDoc.TfIdf) {
			t.Errorf("got %v, wanted %v", tc.doc.TfIdf, tc.expDoc.TfIdf)
		}

		if !approxEq(tc.doc.Norm, tc.expDoc.Norm) {
			t.Errorf("got %.4f, wanted %.4f", tc.doc.Norm, tc.expDoc.Norm)
		}

		if tc.doc.TermFreq != nil {
			t.Errorf("got %v, wanted nil", tc.doc.TermFreq)
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
				TfIdf: TermMap{},
			},
			0.0,
		},
		{
			Document{
				TfIdf: TermMap{
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
