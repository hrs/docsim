package corpus

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/godoc/util"
	"golang.org/x/tools/godoc/vfs"
)

type termMap map[string]float64

type Document struct {
	path     string
	termFreq termMap
	tfIdf    termMap
	norm     float64
}

func ParseDocument(path string, config *Config) (*Document, error) {
	// Ensure that this is a text file
	if !isTextFile(path) {
		return nil, fmt.Errorf("not a text file, skipping: %s", path)
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Parse the contents
	doc, err := NewDocument(file, config)
	if err != nil {
		return nil, err
	}

	doc.path = path
	return doc, nil
}

func NewDocument(rd io.Reader, config *Config) (*Document, error) {
	// Create a scanner from the file
	scanner := bufio.NewScanner(rd)

	// Set the split function for the scanning operation
	scanner.Split(bufio.ScanWords)

	termCount := make(termMap)
	totalWordCount := 0.0

	// Loop over the words, stem them if configured, pass them through the
	// stoplist if configured, and, for each that "should" "count", increment it
	// in the term map.
	for scanner.Scan() {
		token := strings.ToLower(scanner.Text())

		// Split each token on non-alphanumeric characters (except single quotes, to
		// handle contractions)
		for _, word := range strings.FieldsFunc(token, splitToken) {
			// Since we didn't split on single quotes, we need to trim them off now.
			// We'd like "don't" to stay "don't", but "'hello" to become "hello".
			word = strings.Trim(word, "'")

			// Similarly, we need to remove the common "'s" possessive case
			word = strings.TrimSuffix(word, "'s")

			if word != "" && (config.NoStoplist || !config.Stoplist.include(word)) {
				if config.NoStemming {
					termCount[word]++
				} else {
					termCount[stem(word)]++
				}

				totalWordCount++
			}
		}
	}

	// Check for errors in scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Scale the term frequency map according to the total number of terms in the document.
	termFreq := make(termMap)
	for term, count := range termCount {
		termFreq[term] = count / totalWordCount
	}

	return &Document{termFreq: termFreq}, nil
}

func splitToken(r rune) bool {
	return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != ' ' && r != '\''
}

func (doc *Document) normalizeTfIdf(invDocFreq termMap) {
	// Set the TF-IDF weights
	doc.tfIdf = make(termMap)
	for term, weight := range doc.termFreq {
		doc.tfIdf[term] = weight * invDocFreq[term]
	}

	// Now that we've set TF-IDF weights, we can save memory by removing the
	// original weights
	doc.termFreq = nil

	// Calculate and store the document's norm
	doc.norm = doc.calcNorm()
}

func (doc *Document) calcNorm() float64 {
	norm := 0.0
	for _, weight := range doc.tfIdf {
		norm += weight * weight
	}

	return math.Sqrt(norm)
}

func isTextFile(path string) bool {
	mimeType := mime.TypeByExtension(filepath.Ext(path))

	return strings.HasPrefix(mimeType, "text/") ||
		(mimeType == "" && util.IsTextFile(vfs.OS("."), path))
}
