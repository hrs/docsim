package main

import "math"

type Corpus struct {
	Documents  []*Document
	InvDocFreq TermMap
}

func NewCorpus(documents []*Document) *Corpus {
	var docFreq = make(TermMap)
	var invDocFreq = make(TermMap)

	// For each term, in how many documents does it occur?
	for _, doc := range documents {
		for term := range doc.TermFreq {
			docFreq[term]++
		}
	}

	// Invert document frequency and scale
	docCount := float64(len(documents))
	for term := range docFreq {
		invDocFreq[term] = math.Log(docCount / docFreq[term])
	}

	return &Corpus{documents, invDocFreq}
}
