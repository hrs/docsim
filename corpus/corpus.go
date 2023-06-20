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

func ParseCorpus(query *Document, paths []string, config *Config) *Corpus {
	var documents []*Document

	for _, path := range paths {
		err := filepath.WalkDir(path, func(xpath string, xinfo fs.DirEntry, xerr error) error {
			if xerr != nil {
				if errors.Is(xerr, os.ErrNotExist) {
					log.Fatal("no such file or directory: ", path)
				} else {
					panic(xerr)
				}
			}

			// Don't parse directories or symlinks (or the queried file, if so configured)
			if isParsableFile(xinfo, config) && !(config.OmitQuery && sameFile(query.path, xpath)) {
				doc, err := ParseDocument(xpath, config)

				if err != nil {
					if config.Verbose {
						fmt.Fprintln(os.Stderr, err)
					}
				} else {
					documents = append(documents, doc)
				}
			}

			return nil
		})

		if err != nil {
			panic(err)
		}
	}

	return NewCorpus(documents)
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
