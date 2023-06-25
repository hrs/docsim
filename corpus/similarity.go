package corpus

func (corpus *Corpus) SimilarDocuments(query *Document) []score {
	scores := []score{}

	for _, doc := range corpus.documents {
		if doc != query {
			score := score{
				document: doc,
				score:    doc.cosineSimilarity(query),
			}

			scores = append(scores, score)
		}
	}

	return scores
}

func (target *Document) cosineSimilarity(other *Document) float64 {
	dotProd := 0.0

	for term, weight := range target.tfIdf {
		dotProd += (weight * other.tfIdf[term])
	}

	return dotProd / (target.norm * other.norm)
}
