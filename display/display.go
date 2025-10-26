package display

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	counter "bloom.io/github.com/FerDev12/wc-go"
)

type Options struct {
	showLines  bool
	showWords  bool
	showBytes  bool
	showHeader bool
}

type NewOptionsArgs struct {
	ShowLines  bool
	ShowWords  bool
	ShowBytes  bool
	ShowHeader bool
}

func NewOptions(args NewOptionsArgs) Options {
	return Options{
		showLines:  args.ShowLines,
		showWords:  args.ShowWords,
		showBytes:  args.ShowBytes,
		showHeader: args.ShowHeader,
	}
}

func (opts Options) shouldShowAll() bool {
	return opts.showLines == opts.showWords && opts.showWords == opts.showBytes
}

func (opts Options) PrintHeader(w io.Writer) {
	if !opts.showHeader {
		return
	}

	showAll := opts.shouldShowAll()

	if opts.showLines || showAll {
		fmt.Fprintf(w, "lines\t")
	}
	if opts.showWords || showAll {
		fmt.Fprintf(w, "words\t")
	}
	if opts.showBytes || showAll {
		fmt.Fprintf(w, "characters\t")
	}

	fmt.Fprintln(w)
}

func (opts Options) PrintCounts(w io.Writer, c counter.Counts, suffixes ...string) {
	stats := []string{}
	showAll := opts.shouldShowAll()

	if opts.showLines || showAll {
		stats = append(stats, strconv.Itoa(c.Lines))
	}
	if opts.showWords || showAll {
		stats = append(stats, strconv.Itoa(c.Words))
	}
	if opts.showBytes || showAll {
		stats = append(stats, strconv.Itoa(c.Bytes))
	}

	line := strings.Join(stats, "\t") + "\t"
	suffixStr := strings.Join(suffixes, " ")

	fmt.Fprintf(w, "%s", line)

	if suffixStr != "" {
		fmt.Fprintf(w, " %s", suffixStr)
	}

	fmt.Fprintf(w, "\n")
}
