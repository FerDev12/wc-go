package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"
)

type DisplayOptions struct {
	ShowLines  bool
	ShowWords  bool
	ShowBytes  bool
	ShowHeader bool
}

func (opts DisplayOptions) ShouldShowAll() bool {
	return opts.ShowLines == opts.ShowWords && opts.ShowWords == opts.ShowBytes
}

func (opts DisplayOptions) PrintHeader(w io.Writer) {
	if !opts.ShowHeader {
		return
	}

	showAll := opts.ShouldShowAll()

	if opts.ShowLines || showAll {
		fmt.Fprintf(w, "lines\t")
	}

	if opts.ShowWords || showAll {
		fmt.Fprintf(w, "words\t")
	}

	if opts.ShowBytes || showAll {
		fmt.Fprintf(w, "characters\t")
	}

	fmt.Fprintln(w)
}

const TAB_WIDTH = 8
const PADDING = 1
const PAD_CHAR = ' '
const TAB_FLAG = tabwriter.AlignRight

func main() {
	log.SetFlags(0)

	opts := DisplayOptions{}

	flag.BoolVar(&opts.ShowLines, "l", false, "Used to toggle whether or not to show the line count")
	flag.BoolVar(&opts.ShowWords, "w", false, "Used to toggle whether or not to show the word count")
	flag.BoolVar(&opts.ShowBytes, "c", false, "Used to toggle whether or not to show the byte count")
	flag.BoolVar(&opts.ShowHeader, "header", false, "Used to toggle whether or not to show the header")

	flag.Parse()

	// instantiate tabwriter to provide tabular ouptut and define it's behaviour
	wr := tabwriter.NewWriter(os.Stdout, 0, TAB_WIDTH, PADDING, PAD_CHAR, TAB_FLAG)

	filenames := flag.Args()
	didError := false
	totals := Counts{}

	opts.PrintHeader(wr)

	for _, filename := range filenames {
		counts, err := CountFile(filename)

		if err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "counter:", err)
			continue
		}

		counts.Print(wr, opts, filename)

		totals = totals.Add(counts)
	}

	if len(filenames) == 0 {
		GetCounts(os.Stdin).Print(wr, opts)
	} else {
		totals.Print(wr, opts, "total")
	}

	wr.Flush()

	if didError {
		os.Exit(1)
	}
}
