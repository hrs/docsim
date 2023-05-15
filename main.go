package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func makeCorpus(target *Document, paths []string, config Config) *Corpus {
	var documents []*Document

	for _, path := range paths {
		err := filepath.WalkDir(path, func(xpath string, xinfo fs.DirEntry, xerr error) error {
			if xerr != nil {
				panic(xerr)
			}

			if !xinfo.IsDir() && !(config.OmitTarget && sameFile(target.Path, xpath)) {
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
	targetFlag := flag.String("target", "", "path to the file that results should match")
	showScoresFlag := flag.Bool("show-scores", false, "print scores next to file paths")
	bestFirstFlag := flag.Bool("best-first", false, "print best matches first")
	limitFlag := flag.Int("limit", 0, "return at most `limit` results")
	verboseFlag := flag.Bool("verbose", false, "include debugging information and errors")
	omitTargetFlag := flag.Bool("omit-target", false, "don't include the target file itself in search results")
	flag.Parse()

	config := Config{
		ShowScores: *showScoresFlag,
		BestFirst:  *bestFirstFlag,
		OmitTarget: *omitTargetFlag,
		Limit:      *limitFlag,
		Verbose:    *verboseFlag,
	}

	target, _ := NewDocument(*targetFlag)
	corpus := makeCorpus(target, flag.Args(), config)

	printResults(corpus.SimilarDocuments(target), config)
}
