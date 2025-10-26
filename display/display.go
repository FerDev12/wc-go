package display

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	counter "bloom.io/github.com/FerDev12/wc-go"
)

type Options struct {
	ShowLines  bool
	ShowWords  bool
	ShowBytes  bool
	ShowHeader bool
}

func (opts Options) ShouldShowAll() bool {
	return opts.ShowLines == opts.ShowWords && opts.ShowWords == opts.ShowBytes
}

func (opts Options) PrintHeader(w io.Writer) {
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

func (opts Options) PrintCounts(w io.Writer, c counter.Counts, suffixes ...string) {
	stats := []string{}
	showAll := opts.ShouldShowAll()

	if opts.ShowLines || showAll {
		stats = append(stats, strconv.Itoa(c.Lines))
	}
	if opts.ShowWords || showAll {
		stats = append(stats, strconv.Itoa(c.Words))
	}
	if opts.ShowBytes || showAll {
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
