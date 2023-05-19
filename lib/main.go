package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	bestFirstFlag := flag.Bool("best-first", false, "print best matches first")
	limitFlag := flag.Int("limit", 0, "return at most `limit` results")
	noStemmingFlag := flag.Bool("no-stemming", false, "don't perform stemming on words")
	noStoplistFlag := flag.Bool("no-stoplist", false, "don't omit common words by using a stoplist")
	omitQueryFlag := flag.Bool("omit-query", false, "don't include the query file itself in search results")
	queryFlag := flag.String("query", "", "path to the file that results should match")
	showScoresFlag := flag.Bool("show-scores", false, "print scores next to file paths")
	verboseFlag := flag.Bool("verbose", false, "include debugging information and errors")
	flag.Parse()

	config := Config{
		BestFirst:  *bestFirstFlag,
		Limit:      *limitFlag,
		NoStemming: *noStemmingFlag,
		NoStoplist: *noStoplistFlag,
		OmitQuery:  *omitQueryFlag,
		ShowScores: *showScoresFlag,
		Verbose:    *verboseFlag,
	}

	if !config.Verbose {
		// Suppress log timestamps and noisy output
		log.SetFlags(0)
	}

	// If no query file was provided, read from stdin, write to a tempfile, and
	// use that
	queryPath := *queryFlag
	if queryPath == "" {
		reader := bufio.NewReader(os.Stdin)
		data, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatal("Error reading from STDIN:", err)
		}

		f, err := os.CreateTemp("", "docsim-*.txt")
		if err != nil {
			log.Fatal("error creating temporary file:", err)
		}
		defer os.Remove(f.Name())

		_, err = f.Write(data)
		if err != nil {
			log.Fatal("error writing to temporary file:", err)
		}
		f.Close()

		queryPath = f.Name()
	}

	query, err := NewDocument(queryPath, &config)
	if err != nil {
		log.Fatal("error parsing query:", err)
	}

	// If no search paths were provided, search the current directory
	searchPaths := flag.Args()
	if len(searchPaths) == 0 {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("error determining current directory:", err)
		}

		searchPaths = []string{currentDir}
	}

	corpus := ParseCorpus(query, searchPaths, &config)

	printResults(corpus.SimilarDocuments(query), config)
}
