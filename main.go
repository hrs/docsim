package main

import (
	"flag"
	"fmt"
	"os"
)

func makeCorpus(target *Document, paths []string, config Config) *Corpus {
	var documents []*Document

	for _, path := range paths {
		if config.omitTarget && sameFile(target.Path, path) {
			continue
		}

		doc, err := NewDocument(path)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			documents = append(documents, doc)
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
	targetFlag := flag.String("target", "", "path to the file that results should match")
	showScoresFlag := flag.Bool("show-scores", false, "print scores next to file paths")
	bestFirstFlag := flag.Bool("best-first", false, "print best matches first")
	limitFlag := flag.Int("limit", 0, "return at most `limit` results")
	omitTargetFlag := flag.Bool("omit-target", false, "don't include the target file itself in search results")
	flag.Parse()

	config := Config{
		showScores: *showScoresFlag,
		bestFirst:  *bestFirstFlag,
		omitTarget: *omitTargetFlag,
		limit:      *limitFlag,
	}

	target, _ := NewDocument(*targetFlag)
	corpus := makeCorpus(target, flag.Args(), config)

	printResults(corpus.SimilarDocuments(target), config)
}
