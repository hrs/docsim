package corpus

import "github.com/reiver/go-porterstemmer"

var stemCache = make(map[string]string)

func stem(word string) string {
	cachedStem, ok := stemCache[word]

	if ok {
		return cachedStem
	}

	stem := string(porterstemmer.StemWithoutLowerCasing([]rune(word)))
	stemCache[word] = stem
	return stem
}
