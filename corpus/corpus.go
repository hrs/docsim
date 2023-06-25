package corpus

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
)

type Corpus struct {
	documents  []*Document
	invDocFreq termMap
}

func NewCorpus(documents []*Document) *Corpus {
	var docFreq = make(termMap)
	var invDocFreq = make(termMap)

	// For each term, in how many documents does it occur?
	for _, doc := range documents {
		for term := range doc.termFreq {
			docFreq[term]++
		}
	}

	// Invert document frequency and scale
	docCount := float64(len(documents))
	for term := range docFreq {
		invDocFreq[term] = math.Log(docCount / docFreq[term])
	}

	// Assign TF-IDF weights to every document in the corpus
	for _, doc := range documents {
		doc.normalizeTfIdf(invDocFreq)
	}

	return &Corpus{documents, invDocFreq}
}

func ParseCorpus(query *Document, searchPaths []string, config *Config) *Corpus {
	documents := []*Document{query}

	for _, searchPath := range searchPaths {
		documents = append(documents, parseDocuments(query, searchPath, config)...)
	}

	return NewCorpus(documents)
}

func parseDocuments(query *Document, searchPath string, config *Config) []*Document {
	var documents []*Document

	for _, path := range findParsableFiles(searchPath, config) {
		doc, err := ParseDocument(path, config)

		if err != nil {
			if config.Verbose {
				fmt.Fprintln(os.Stderr, err)
			}
		} else {
			documents = append(documents, doc)
		}
	}

	return documents
}

func findParsableFiles(searchPath string, config *Config) []string {
	filePaths := []string{}

	err := filepath.WalkDir(searchPath, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Fatal("no such file or directory: ", searchPath)
			} else {
				panic(err)
			}
		}

		// Don't include directories (or symlinks, unless so configured).
		if isParsableFile(info, config) {
			filePaths = append(filePaths, path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return filePaths
}

func isParsableFile(info fs.DirEntry, config *Config) bool {
	// Return true if this is either a regular file OR if it's a symlink that the
	// user's chosen to include
	return info.Type().IsRegular() ||
		(config.FollowSymlinks && (info.Type()&os.ModeSymlink != 0))
}
