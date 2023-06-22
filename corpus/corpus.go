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
	var documents []*Document

	for _, filepath := range findParsableFiles(searchPaths, config) {
		// Don't include the queried file, if so configured.
		if !(config.OmitQuery && sameFile(query.path, filepath)) {
			doc, err := ParseDocument(filepath, config)

			if err != nil {
				if config.Verbose {
					fmt.Fprintln(os.Stderr, err)
				}
			} else {
				documents = append(documents, doc)
			}
		}
	}

	return NewCorpus(documents)
}

func findParsableFiles(searchPaths []string, config *Config) []string {
	filePaths := []string{}

	for _, path := range searchPaths {
		err := filepath.WalkDir(path, func(xpath string, xinfo fs.DirEntry, xerr error) error {
			if xerr != nil {
				if errors.Is(xerr, os.ErrNotExist) {
					log.Fatal("no such file or directory: ", path)
				} else {
					panic(xerr)
				}
			}

			// Don't include directories (or symlinks, unless so configured).
			if isParsableFile(xinfo, config) {
				filePaths = append(filePaths, xpath)
			}

			return nil
		})

		if err != nil {
			panic(err)
		}
	}

	return filePaths
}

func isParsableFile(info fs.DirEntry, config *Config) bool {
	// Return true if this is either a regular file OR if it's a symlink that the
	// user's chosen to include
	return info.Type().IsRegular() ||
		(config.FollowSymlinks && (info.Type()&os.ModeSymlink != 0))
}

func sameFile(a, b string) bool {
	aFileInfo, err := os.Stat(a)
	if err != nil {
		return false
	}

	bFileInfo, err := os.Stat(b)
	if err != nil {
		return false
	}

	return os.SameFile(aFileInfo, bFileInfo)
}
