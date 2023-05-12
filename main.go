package main

import (
	"flag"
	"fmt"
)

func main() {
	targetFlag := flag.String("target", "", "path to the file that results should match")
	flag.Parse()

	docs := flag.Args()

	fmt.Println("target: ", *targetFlag)
	fmt.Println("corpus:", docs)
}
