package corpus

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

	docA := Document{tfIdf: termMap{0: 0.3013, 1: 0.2628}}
	docA.norm = docA.calcNorm()

	docB := Document{tfIdf: termMap{2: 0.1577, 3: 0.7796, 4: 0.1577}}
	docB.norm = docB.calcNorm()

	docC := Document{tfIdf: termMap{0: 0.2260, 3: 0.6496}}
	docC.norm = docC.calcNorm()

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
