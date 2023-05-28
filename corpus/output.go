package corpus

import (
	"fmt"
	"sort"
)

type score struct {
	query    *Document
	document *Document
	score    float64
}

func PrintResults(scores []score, config Config) {
	// Sort results by score, worst matches first
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	if config.Limit > 0 && len(scores) > config.Limit {
		scores = scores[0:config.Limit]
	}

	if !config.BestFirst {
		tmp := make([]score, len(scores))
		for i, score := range scores {
			tmp[len(scores)-i-1] = score
		}
		scores = tmp
	}

	for _, score := range scores {
		if config.ShowScores {
			fmt.Printf("%.4f\t%s\n", score.score, score.document.path)
		} else {
			fmt.Println(score.document.path)
		}
	}
}
