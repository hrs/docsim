package main

type Corpus struct {
	Documents []*Document
	Terms     TermMap
}

func NewCorpus(documents []*Document) *Corpus {
	var terms = make(TermMap)

	for _, doc := range documents {
		for term, count := range doc.Terms {
			terms[term] += count
		}
	}

	return &Corpus{documents, terms}
}
