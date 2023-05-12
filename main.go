package main

import (
	"flag"
	"fmt"
	"log"
)

func makeCorpus(paths []string) *Corpus {
	var documents []*Document

	for _, path := range paths {
		doc, err := NewDocument(path)

		if err != nil {
			log.Println(err)
		} else {
			documents = append(documents, doc)
		}
	}

	return NewCorpus(documents)
}

func main() {
	targetFlag := flag.String("target", "", "path to the file that results should match")
	flag.Parse()

	corpus := makeCorpus(flag.Args())

	fmt.Println("target:", *targetFlag)
	fmt.Println("corpus:", corpus)
}
