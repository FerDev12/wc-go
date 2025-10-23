package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	totalLines := 0
	totalWords := 0
	totalBytes := 0

	filenames := os.Args[1:]
	didError := false

	for _, filename := range filenames {
		res, err := CountFile(filename)

		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}

		fmt.Println(filename, "-", res)

		totalLines += res.Lines
		totalWords += res.Words
		totalBytes += res.Bytes
	}

	if len(filenames) == 0 {
		res := CountWords(os.Stdin)
		fmt.Println(res)
	} else {
		fmt.Println("total -", CountsResult{
			Lines: totalLines,
			Words: totalWords,
			Bytes: totalBytes,
		})
	}

	if didError {
		os.Exit(1)
	}
}
