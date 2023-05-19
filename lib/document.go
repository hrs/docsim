package main

import (
	"bufio"
	"fmt"
	"math"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/godoc/util"
	"golang.org/x/tools/godoc/vfs"
)

type TermMap map[string]float64

type Document struct {
	Path     string
	TermFreq TermMap
	TfIdf    TermMap
	Norm     float64
}

func NewDocument(path string, config *Config) (*Document, error) {
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

	// Create a scanner from the file
	scanner := bufio.NewScanner(file)

	// Set the split function for the scanning operation
	scanner.Split(bufio.ScanWords)

	// Initialize the words slice
	termCount := make(TermMap)
	totalWordCount := 0.0

	// Loop over the words and append each to the words slice
	for scanner.Scan() {
		token := strings.ToLower(scanner.Text())

		// Split each token on non-alphanumeric characters (except single qutoes, to
		// handle contractions)
		for _, word := range strings.FieldsFunc(token, splitToken) {
			// Since we didn't split on single quotes, we need to trim them off now.
			// We'd like "don't" to stay "don't", but "'hello" to become "hello".
			word = strings.Trim(word, "'")

			// Similarly, we need to remove the common "'s" possessive case
			word = strings.TrimSuffix(word, "'s")

			if word != "" {
				if config.NoStoplist || !InStoplist(word) {
					if config.NoStemming {
						termCount[word]++
					} else {
						termCount[stem(word)]++
					}

					totalWordCount++
				}
			}
		}
	}

	// Check for errors in scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Build the term frequency map
	termFreq := make(TermMap)
	for term, count := range termCount {
		termFreq[term] = count / totalWordCount
	}

	return &Document{Path: path, TermFreq: termFreq}, nil
}

func splitToken(r rune) bool {
	return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9') && r != ' ' && r != '\''
}

func (doc *Document) NormalizeTfIdf(invDocFreq TermMap) {
	// Set the TF-IDF weights
	doc.TfIdf = make(TermMap)
	for term, weight := range doc.TermFreq {
		doc.TfIdf[term] = weight * invDocFreq[term]
	}

	// Now that we've set TF-IDF weights, we can save memory by removing the
	// original weights
	doc.TermFreq = nil

	// Calculate and store the document's norm
	doc.Norm = doc.calcNorm()
}

func (doc *Document) calcNorm() float64 {
	norm := 0.0
	for _, weight := range doc.TfIdf {
		norm += weight * weight
	}

	return math.Sqrt(norm)
}

func isTextFile(path string) bool {
	mimeType := mime.TypeByExtension(filepath.Ext(path))

	return strings.HasPrefix(mimeType, "text/") ||
		(mimeType == "" && util.IsTextFile(vfs.OS("."), path))
}
