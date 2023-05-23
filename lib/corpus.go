package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
)

type Corpus struct {
	Documents  []*Document
	InvDocFreq TermMap
}

func NewCorpus(documents []*Document) *Corpus {
	var docFreq = make(TermMap)
	var invDocFreq = make(TermMap)

	// For each term, in how many documents does it occur?
	for _, doc := range documents {
		for term := range doc.TermFreq {
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
		doc.NormalizeTfIdf(invDocFreq)
	}

	return &Corpus{documents, invDocFreq}
}

func ParseCorpus(query *Document, paths []string, config *Config) *Corpus {
	var documents []*Document

	for _, path := range paths {
		err := filepath.WalkDir(path, func(xpath string, xinfo fs.DirEntry, xerr error) error {
			if xerr != nil {
				panic(xerr)
			}

			// Don't parse directories or symlinks (or the queried file, if so configured)
			if xinfo.Type().IsRegular() && !(config.OmitQuery && sameFile(query.Path, xpath)) {
				doc, err := NewDocument(xpath, config)

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
