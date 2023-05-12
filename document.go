package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Document struct {
	Path  string
	Words []string
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-z0-9 ']+`)

func NewDocument(path string) (*Document, error) {
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
	words := []string{}

	// Loop over the words and append each to the words slice
	for scanner.Scan() {
		token := strings.ToLower(scanner.Text())

		// Split each token on non-alphanumeric characters (except single qutoes, to
		// handle contractions)
		for _, word := range nonAlphanumericRegex.Split(token, -1) {
			// Since we didn't split on single quotes, we need to trim them off now.
			// We'd like "don't" to stay "don't", but "'hello" to become "hello".
			word = strings.Trim(word, "'")

			// Similarly, we need to remove the common "'s" possessive case
			word = strings.TrimSuffix(word, "'s")

			if word != "" && !inStoplist(word) {
				words = append(words, stem(word))
			}
		}
	}

	// Check for errors in scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Document{path, words}, nil
}
