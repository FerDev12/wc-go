package display

import (
	"fmt"
	"io"
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

func (opts Options) ShouldShowLines() bool {
	return opts.args.ShowLines || opts.shouldShowAll()
}

func (opts Options) ShouldShowWords() bool {
	return opts.args.ShowWords || opts.shouldShowAll()
}

func (opts Options) ShouldShowBytes() bool {
	return opts.args.ShowBytes || opts.shouldShowAll()
}

func (opts Options) PrintHeader(w io.Writer) {
	if !opts.args.ShowHeader {
		return
	}

	if opts.ShouldShowLines() {
		fmt.Fprintf(w, "lines\t")
	}
	if opts.ShouldShowWords() {
		fmt.Fprintf(w, "words\t")
	}
	if opts.ShouldShowBytes() {
		fmt.Fprintf(w, "characters\t")
	}

	fmt.Fprintln(w)
}
