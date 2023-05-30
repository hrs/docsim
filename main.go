package main

import (
	"flag"
	"log"
	"os"

	"github.com/hrs/docsim/corpus"
)

func stoplist(flag string) *corpus.Stoplist {
	if flag == "" {
		return corpus.DefaultStoplist
	} else {
		var err error
		stoplist, err := corpus.ParseStoplist(flag)
		if err != nil {
			log.Fatal("Error reading custom stoplist:", err)
		}
		return stoplist
	}
}

func queryDoc(path string, config *corpus.Config) (*corpus.Document, error) {
	if path == "" {
		return corpus.NewDocument(os.Stdin, config)
	}

	return corpus.ParseDocument(path, config)
}

func main() {
	bestFirstFlag := flag.Bool("best-first", false, "print best matches first")
	followSymlinksFlag := flag.Bool("follow-symlinks", false, "included symlinked files in results")
	limitFlag := flag.Int("limit", 0, "return at most `limit` results")
	noStemmingFlag := flag.Bool("no-stemming", false, "don't perform stemming on words")
	noStoplistFlag := flag.Bool("no-stoplist", false, "don't omit common words by using a stoplist")
	omitQueryFlag := flag.Bool("omit-query", false, "don't include the query file itself in search results")
	queryFlag := flag.String("query", "", "path to the file that results should match")
	showScoresFlag := flag.Bool("show-scores", false, "print scores next to file paths")
	stoplistFlag := flag.String("stoplist", "", "path to a file of words to be ignored")
	verboseFlag := flag.Bool("verbose", false, "include debugging information and errors")
	flag.Parse()

	config := corpus.Config{
		BestFirst:      *bestFirstFlag,
		FollowSymlinks: *followSymlinksFlag,
		Limit:          *limitFlag,
		NoStemming:     *noStemmingFlag,
		NoStoplist:     *noStoplistFlag,
		OmitQuery:      *omitQueryFlag,
		ShowScores:     *showScoresFlag,
		Stoplist:       stoplist(*stoplistFlag),
		Verbose:        *verboseFlag,
	}

	if !config.Verbose {
		// Suppress log timestamps and noisy output
		log.SetFlags(0)
	}

	// If no query file was provided, read from stdin
	query, err := queryDoc(*queryFlag, &config)
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

	c := corpus.ParseCorpus(query, searchPaths, &config)

	corpus.PrintResults(c.SimilarDocuments(query), config)
}
