package main

import (
	"bufio"
	"os"
	"strings"
)

type Document struct {
	Path  string
	Words []string
}

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
		words = append(words, strings.ToLower(scanner.Text()))
	}

	// Check for errors in scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &Document{path, words}, nil
}
