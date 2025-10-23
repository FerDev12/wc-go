package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(0)

	totalWords := 0
	filenames := os.Args[1:]
	didError := false
	start := time.Now()

	for _, filename := range filenames {
		wordCount, err := CountWordsForFile(filename)

		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter: ", err)
			continue
		}

		fmt.Println(fmt.Sprintf("%d %s", wordCount, filename))
		totalWords += wordCount
	}

	if len(filenames) == 0 {
		wordCount := CountWords(os.Stdin)
		fmt.Println(wordCount)
	} else {
		fmt.Println(fmt.Sprintf("%d total %v seconds", totalWords, time.Since(start)))
	}

	if didError {
		os.Exit(1)
	}
}
