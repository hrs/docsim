package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func makeCorpus(query *Document, paths []string, config Config) *Corpus {
	var documents []*Document

	for _, path := range paths {
		err := filepath.WalkDir(path, func(xpath string, xinfo fs.DirEntry, xerr error) error {
			if xerr != nil {
				panic(xerr)
			}

			if !xinfo.IsDir() && !(config.OmitQuery && sameFile(query.Path, xpath)) {
				doc, err := NewDocument(xpath)

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

func main() {
	queryFlag := flag.String("query", "", "path to the file that results should match")
	showScoresFlag := flag.Bool("show-scores", false, "print scores next to file paths")
	bestFirstFlag := flag.Bool("best-first", false, "print best matches first")
	limitFlag := flag.Int("limit", 0, "return at most `limit` results")
	verboseFlag := flag.Bool("verbose", false, "include debugging information and errors")
	omitQueryFlag := flag.Bool("omit-query", false, "don't include the query file itself in search results")
	flag.Parse()

	config := Config{
		ShowScores: *showScoresFlag,
		BestFirst:  *bestFirstFlag,
		OmitQuery:  *omitQueryFlag,
		Limit:      *limitFlag,
		Verbose:    *verboseFlag,
	}

	query, _ := NewDocument(*queryFlag)
	corpus := makeCorpus(query, flag.Args(), config)

	printResults(corpus.SimilarDocuments(query), config)
}
