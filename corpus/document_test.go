package corpus

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseTokens(t *testing.T) {
	tests := []struct {
		text string
		exp  []string
	}{
		{
			"naïve née señor",
			[]string{"naïve", "née", "señor"},
		},
		{
			"1337 hAx0r",
			[]string{"1337", "hax0r"},
		},
		{
			"examples: isn't 'isn't' wasn’t 'wasn’t' ‘won't’ ‘won't’ ‘shan’t’ ‘shan’t’",
			[]string{"examples", "isn't", "isn't", "wasn't", "wasn't", "won't", "won't", "shan't", "shan't"},
		},
	}

	for _, tc := range tests {
		got, err := parseTokens(strings.NewReader(tc.text))
		if err != nil {
			t.Errorf("got unexpected error %v", err)
		}

		if !reflect.DeepEqual(got, tc.exp) {
			t.Errorf("got %#v, wanted %#v", got, tc.exp)
		}
	}
}

func TestNewDocument(t *testing.T) {
	sampleText := "It had two positions, and scrawled in pencil on the metal switch body were the words 'magic' and 'more magic'."

	tests := []struct {
		config Config
		expMap termMap
	}{
		{
			Config{Stoplist: DefaultStoplist},
			termMap{
				"bodi":   0.1250,
				"magic":  0.2500,
				"metal":  0.1250,
				"pencil": 0.1250,
				"posit":  0.1250,
				"scrawl": 0.1250,
				"switch": 0.1250,
			},
		},
		{
			Config{
				Stoplist: newStoplist(
					[]string{
						"and",
						"in",
						"on",
						"the",
						"were",
					},
				),
			},
			termMap{
				"bodi":   0.0769,
				"had":    0.0769,
				"it":     0.0769,
				"magic":  0.1538,
				"metal":  0.0769,
				"more":   0.0769,
				"pencil": 0.0769,
				"posit":  0.0769,
				"scrawl": 0.0769,
				"switch": 0.0769,
				"two":    0.0769,
				"word":   0.0769,
			},
		},
		{
			Config{
				NoStoplist: true,
			},
			termMap{
				"and":    0.1000,
				"bodi":   0.0500,
				"had":    0.0500,
				"in":     0.0500,
				"it":     0.0500,
				"magic":  0.1000,
				"metal":  0.0500,
				"more":   0.0500,
				"on":     0.0500,
				"pencil": 0.0500,
				"posit":  0.0500,
				"scrawl": 0.0500,
				"switch": 0.0500,
				"the":    0.1000,
				"two":    0.0500,
				"were":   0.0500,
				"word":   0.0500,
			},
		},
		{
			Config{
				NoStemming: true,
				Stoplist:   DefaultStoplist,
			},
			termMap{
				"body":      0.1250,
				"magic":     0.2500,
				"metal":     0.1250,
				"pencil":    0.1250,
				"positions": 0.1250,
				"scrawled":  0.1250,
				"switch":    0.1250,
			},
		},
	}

	for _, tc := range tests {
		got, err := NewDocument(strings.NewReader(sampleText), &tc.config)
		if err != nil {
			t.Errorf("got unexpected error %v", err)
		}

		for gotTerm := range got.termFreq {
			_, expKey := tc.expMap[gotTerm]
			if !expKey {
				t.Errorf("parsed unexpected term '%s'", gotTerm)
			}
		}

		for expTerm, expFreq := range tc.expMap {
			gotFreq, ok := got.termFreq[expTerm]
			if !ok {
				t.Errorf("found unexpected term '%s' in termFreq", expTerm)
			}

			if !approxEq(gotFreq, expFreq) {
				t.Errorf("for term '%s' got %.4f, wanted %.4f", expTerm, gotFreq, expFreq)
			}
		}
	}
}

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
