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
const PADDING = 1
const PAD_CHAR = ' '
const TAB_FLAG = tabwriter.AlignRight

type FilesCountResult struct {
	counts   counter.Counts
	filename string
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

	wg := sync.WaitGroup{}
	wg.Add(len(filenames))

	ch := make(chan FilesCountResult)

	for _, filename := range filenames {
		go func() {
			defer wg.Done()

			counts, err := counter.CountFile(filename)
			if err != nil {
				didError = true
				fmt.Fprintln(os.Stderr, "counter:", err)
				return
			}

			ch <- FilesCountResult{
				counts:   counts,
				filename: filename,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
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
