package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type DisplayOptions struct {
	ShowLines bool
	ShowWords bool
	ShowBytes bool
}

func (opts DisplayOptions) ShouldShowAll() bool {
	return opts.ShowLines == opts.ShowWords && opts.ShowWords == opts.ShowBytes
}

func main() {
	log.SetFlags(0)

	opts := DisplayOptions{}

	flag.BoolVar(&opts.ShowLines, "l", false, "Used to toggle whether or not to show the line count")
	flag.BoolVar(&opts.ShowWords, "w", false, "Used to toggle whether or not to show the word count")
	flag.BoolVar(&opts.ShowBytes, "b", false, "Used to toggle whether or not to show the byte count")

	flag.Parse()

	filenames := flag.Args()
	didError := false
	totals := Counts{}

	for _, filename := range filenames {
		counts, err := CountFile(filename)

		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}

		counts.Print(os.Stdout, opts, filename)

		totals = totals.Add(counts)
	}

	if len(filenames) == 0 {
		GetCounts(os.Stdin).Print(os.Stdout, opts)
	} else {
		totals.Print(os.Stdout, opts, "total")
	}

	if didError {
		os.Exit(1)
	}
}
