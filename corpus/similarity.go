package corpus

import "math"

func (corpus *Corpus) SimilarDocuments(query *Document) []score {
	// Normalize query document to set TF-IDF weights per the corpus
	query.normalizeTfIdf(corpus.invDocFreq)

	scores := make([]score, len(corpus.documents))
	for i, doc := range corpus.documents {
		score := score{
			query:    query,
			document: doc,
			score:    doc.cosineSimilarity(query),
		}

		scores[i] = score
	}

	return scores
}

func (target *Document) cosineSimilarity(other *Document) float64 {
	dotProd := 0.0

	for term, weight := range target.tfIdf {
		dotProd += (weight * other.tfIdf[term])
	}

	sim := dotProd / (target.norm * other.norm)
	if math.IsNaN(sim) {
		return 0.0
	}

	return sim
}
