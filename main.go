package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hrs/docsim/corpus"
)

var version = "0.1.4"

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

func failWithUsageMessage() {
	executable := os.Args[0]

	log.Fatalf(
		"usage:\n"+
			"\n"+
			"    %s [OPTION...] QUERY [PATH...]\n"+
			"    %s [OPTION...] --file PATH [PATH...]\n"+
			"    command | %s [OPTION...] --stdin [PATH...]\n"+
			"\n"+
			"For more information try --help, or check the manual.",
		executable,
		executable,
		executable,
	)
}

func main() {
	var config corpus.Config

	flag.BoolVar(&config.BestFirst, "best-first", false, "print best matches first")
	flag.BoolVar(&config.FollowSymlinks, "follow-symlinks", false, "included symlinked files in results")
	flag.IntVar(&config.Limit, "limit", 0, "return at most `limit` results")
	flag.BoolVar(&config.NoStemming, "no-stemming", false, "don't perform stemming on words")
	flag.BoolVar(&config.NoStoplist, "no-stoplist", false, "don't omit common words by using a stoplist")
	flag.BoolVar(&config.OmitQuery, "omit-query", false, "don't include the query file itself in search results")
	flag.BoolVar(&config.ShowScores, "show-scores", false, "print scores next to file paths")
	flag.BoolVar(&config.Verbose, "verbose", false, "include debugging information and errors")
	stdinFlag := flag.Bool("stdin", false, "read query from STDIN instead of from a positional string arugment")
	fileFlag := flag.String("file", "", "path to the file that results should match")
	stoplistFlag := flag.String("stoplist", "", "path to a file of words to be ignored")
	versionFlag := flag.Bool("version", false, "print the current version and exit")
	flag.Parse()

	config.Stoplist = stoplist(*stoplistFlag)

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if !config.Verbose {
		// Suppress log timestamps and noisy output
		log.SetFlags(0)
	}

	positionalArgs := flag.Args()

	// Construct a query from either:
	//
	// - `--stdin`,
	// - `--file`, or
	// - the first positional argument
	var query *corpus.Document
	var err error
	if *stdinFlag {
		// Don't try to read from both `--stdin` and a `--file`
		if *fileFlag != "" {
			log.Println("error: can't read query from both --stdin and a --file")
			failWithUsageMessage()
		}

		query, err = corpus.NewDocument(os.Stdin, &config)
		if err != nil {
			log.Fatal("error parsing STDIN:", err)
		}
	} else if *fileFlag != "" {
		query, err = corpus.ParseDocument(*fileFlag, &config)
		if err != nil {
			log.Fatal("error parsing file:", err)
		}
	} else {
		if len(positionalArgs) == 0 {
			log.Println("error: no query found")
			failWithUsageMessage()
		}

		query, err = corpus.NewDocument(strings.NewReader(positionalArgs[0]), &config)
		if err != nil {
			log.Fatal("error parsing query:", err)
		}

		// Remove the query term from the list of paths to be searched
		positionalArgs = positionalArgs[1:]
	}

	// If no search paths were provided, search the current directory
	if len(positionalArgs) == 0 {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal("error determining current directory:", err)
		}

		positionalArgs = []string{currentDir}
	}

	c := corpus.ParseCorpus(query, positionalArgs, &config)

	corpus.PrintResults(c.SimilarDocuments(query), config)
}