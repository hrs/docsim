package main

import (
	"math"
	"testing"
)

func approxEq(a, b float64) bool {
	maxDelta := 0.0001

	return math.Abs(a-b) < maxDelta
}

func TestCosineSimilarity(t *testing.T) {
	type cosTest struct {
		query  *Document
		target *Document
		sim    float64
	}

	docA := Document{TfIdf: TermMap{"foo": 0.3013, "bar": 0.2628}}
	docA.Norm = docA.calcNorm()

	docB := Document{TfIdf: TermMap{"baz": 0.1577, "quux": 0.7796, "xyzzy": 0.1577}}
	docB.Norm = docB.calcNorm()

	docC := Document{TfIdf: TermMap{"foo": 0.2260, "quux": 0.6496}}
	docC.Norm = docC.calcNorm()

	cosTests := []cosTest{
		{&docA, &docA, 1.0},
		{&docA, &docB, 0.0},
		{&docA, &docC, 0.2476},
	}

	for _, tc := range cosTests {
		sim := tc.query.cosineSimilarity(tc.target)

		if !approxEq(sim, tc.sim) {
			t.Errorf("got %.4f, wanted %.4f", sim, tc.sim)
		}
	}
}
