package corpus

type termMap map[int]float64

var termIDs = make(map[string]int)

func termID(term string) int {
	cachedID, ok := termIDs[term]

	if ok {
		return cachedID
	}

	id := len(termIDs)
	termIDs[term] = id
	return id
}
