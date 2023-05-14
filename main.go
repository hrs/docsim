package main

import (
	"flag"
	"fmt"
	"os"
)

func makeCorpus(paths []string) *Corpus {
	var documents []*Document

	for _, path := range paths {
		doc, err := NewDocument(path)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			documents = append(documents, doc)
		}
	}

	return NewCorpus(documents)
}

func main() {
	targetFlag := flag.String("target", "", "path to the file that results should match")
	showScoresFlag := flag.Bool("show-scores", false, "print scores next to file paths")
	bestFirstFlag := flag.Bool("best-first", false, "print best matches first")
	limitFlag := flag.Int("limit", 0, "return at most `limit` results")
	flag.Parse()

	corpus := makeCorpus(flag.Args())
	target, _ := NewDocument(*targetFlag)

	printResults(
		corpus.SimilarDocuments(target),
		Config{
			showScores: *showScoresFlag,
			bestFirst:  *bestFirstFlag,
			limit:      *limitFlag,
		},
	)
}
