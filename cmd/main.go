package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"text/tabwriter"

	counter "bloom.io/github.com/FerDev12/wc-go"
	"bloom.io/github.com/FerDev12/wc-go/display"
)

const TAB_WIDTH = 8
const PADDING = 4
const PAD_CHAR = ' '
const TAB_FLAG = tabwriter.AlignRight

type FilesCountResult struct {
	counts   counter.Counts
	filename string
	err      error
	idx      int
}

func main() {
	log.SetFlags(0)

	displayOptionsArgs := display.NewOptionsArgs{}

	flag.BoolVar(&displayOptionsArgs.ShowLines, "l", false, "Used to toggle whether or not to show the line count")
	flag.BoolVar(&displayOptionsArgs.ShowBytes, "w", false, "Used to toggle whether or not to show the word count")
	flag.BoolVar(&displayOptionsArgs.ShowBytes, "c", false, "Used to toggle whether or not to show the byte count")
	flag.BoolVar(&displayOptionsArgs.ShowHeader, "header", false, "Used to toggle whether or not to show the header")

	flag.Parse()

	opts := display.NewOptions(displayOptionsArgs)

	// instantiate tabwriter to provide tabular ouptut and define it's behaviour
	wr := tabwriter.NewWriter(os.Stdout, 0, TAB_WIDTH, PADDING, PAD_CHAR, TAB_FLAG)

	filenames := flag.Args()
	didError := false
	totals := counter.Counts{}

	opts.PrintHeader(wr)

	ch := CountFiles(filenames)

	results := make([]FilesCountResult, len(filenames))

	for res := range ch {
		results[res.idx] = res
	}

	for _, res := range results {
		if res.err != nil {
			didError = true
			fmt.Fprintln(os.Stderr, "wc-go:", res.err)
			continue
		}
		totals = totals.Add(res.counts)
		res.counts.Print(wr, opts, res.filename)
	}

	if len(filenames) == 0 {
		counts := counter.GetCounts(os.Stdin)
		counts.Print(wr, opts)
	} else {
		totals.Print(wr, opts, "total")
	}

	wr.Flush()

	if didError {
		os.Exit(1)
	}
}

func CountFiles(filenames []string) <-chan FilesCountResult {
	ch := make(chan FilesCountResult)

	wg := sync.WaitGroup{}
	wg.Add(len(filenames))

	for i, filename := range filenames {
		go func() {
			defer wg.Done()
			counts, err := counter.CountFile(filename)
			ch <- FilesCountResult{
				filename: filename,
				counts:   counts,
				err:      err,
				idx:      i,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
