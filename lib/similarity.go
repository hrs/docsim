package main

func (corpus *Corpus) SimilarDocuments(query *Document) []Score {
	// Normalize query document to set TF-IDF weights per the corpus
	query.NormalizeTfIdf(corpus.InvDocFreq)

	scores := make([]Score, len(corpus.Documents))
	for i, doc := range corpus.Documents {
		score := Score{
			Query:    query,
			Document: doc,
			Score:    doc.cosineSimilarity(query),
		}

		scores[i] = score
	}

	return scores
}

func (target *Document) cosineSimilarity(other *Document) float64 {
	dotProd := 0.0

	for term, weight := range target.TfIdf {
		dotProd += (weight * other.TfIdf[term])
	}

	return dotProd / (target.Norm * other.Norm)
}
