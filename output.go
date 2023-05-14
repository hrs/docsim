package main

import (
	"fmt"
	"sort"
)

type Score struct {
	Query    *Document
	Document *Document
	Score    float64
}

func printResults(scores []Score, config Config) {
	// Sort results by score, worst matches first
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if config.Limit > 0 && len(scores) > config.Limit {
		scores = scores[0:config.Limit]
	}

	if !config.BestFirst {
		tmp := make([]Score, len(scores))
		for i, score := range scores {
			tmp[len(scores)-i-1] = score
		}
		scores = tmp
	}

	for _, score := range scores {
		if config.ShowScores {
			fmt.Printf("%.4f\t%s\n", score.Score, score.Document.Path)
		} else {
			fmt.Println(score.Document.Path)
		}
	}
}
