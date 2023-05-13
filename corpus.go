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

	// Assign TF-IDF weights to every document in the corpus
	for _, doc := range documents {
		doc.NormalizeTfIdf(invDocFreq)
	}

	return &Corpus{documents, invDocFreq}
}

func (corpus *Corpus) SimilarDocuments(query *Document) map[string]float64 {
	// Normalize query document to set TF-IDF weights per the corpus
	query.NormalizeTfIdf(corpus.InvDocFreq)

	similarities := make(map[string]float64)
	for _, doc := range corpus.Documents {
		similarities[doc.Path] = doc.cosineSimilarity(query)
	}

	return similarities
}
