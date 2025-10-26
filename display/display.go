package display

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	counter "bloom.io/github.com/FerDev12/wc-go"
)

type Options struct {
	args NewOptionsArgs
}

type NewOptionsArgs struct {
	ShowLines  bool
	ShowWords  bool
	ShowBytes  bool
	ShowHeader bool
}

func NewOptions(args NewOptionsArgs) Options {
	return Options{
		args: args,
	}
}

func (opts Options) shouldShowAll() bool {
	return opts.args.ShowLines == opts.args.ShowWords && opts.args.ShowWords == opts.args.ShowBytes
}

func (opts Options) PrintHeader(w io.Writer) {
	if !opts.args.ShowHeader {
		return
	}

	showAll := opts.shouldShowAll()

	if opts.args.ShowLines || showAll {
		fmt.Fprintf(w, "lines\t")
	}
	if opts.args.ShowWords || showAll {
		fmt.Fprintf(w, "words\t")
	}
	if opts.args.ShowBytes || showAll {
		fmt.Fprintf(w, "characters\t")
	}

	fmt.Fprintln(w)
}

func (opts Options) PrintCounts(w io.Writer, c counter.Counts, suffixes ...string) {
	stats := []string{}
	showAll := opts.shouldShowAll()

	if opts.args.ShowLines || showAll {
		stats = append(stats, strconv.Itoa(c.Lines))
	}
	if opts.args.ShowWords || showAll {
		stats = append(stats, strconv.Itoa(c.Words))
	}
	if opts.args.ShowBytes || showAll {
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
