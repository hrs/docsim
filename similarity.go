package main

func (target *Document) cosineSimilarity(other *Document) float64 {
	dotProd := 0.0

	for term, weight := range target.TfIdf {
		dotProd += (weight * other.TfIdf[term])
	}

	return dotProd / (target.Norm * other.Norm)
}
