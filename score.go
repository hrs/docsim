package main

import "fmt"

type Score struct {
	Query    *Document
	Document *Document
	Score    float64
}

func (score Score) String() string {
	return fmt.Sprintf("%.4f\t%s", score.Score, score.Document.Path)
}
