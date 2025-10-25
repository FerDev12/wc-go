package main_test

import (
	"bytes"
	"strings"
	"testing"

	counter "bloom.io/word-counter"
)

func TestCountWords(t *testing.T) {
	var testCases = []struct {
		name  string
		input string
		wants int
	}{
		{name: "five words", input: "one two three four five", wants: 5},
		{name: "empty string", input: "", wants: 0},
		{name: "single space", input: " ", wants: 0},
		{name: "new line", input: "one\ntwo", wants: 2},
		{name: "multiple spaces", input: "one   two", wants: 2},
		{name: "prefixed multiple spaces", input: "   one two", wants: 2},
		{name: "suffixed multiple spaces", input: "one two   ", wants: 2},
		{name: "tab characters", input: "	one two		three", wants: 3},
		{name: "utf8 characters", input: "one two three four five six", wants: 6},
		{name: "unicode characters", input: "one two thrРee four five", wants: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			got := counter.GetCounts(reader).Words
			if got != tc.wants {
				t.Errorf("got %d, wants %d", got, tc.wants)
			}
		})
	}

}

func TestCountLines(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wants int
	}{
		{name: "simple five words, 1 new Line", input: "one two three four five\n", wants: 1},
		{name: "empty file", input: "", wants: 0},
		{name: "no new lines", input: "one two three four five", wants: 0},
		{name: "no new line at end", input: "one two three four five\n six", wants: 1},
		{name: "multi newline string", input: "\n\n\n\n", wants: 4},
		{name: "multi word and newline string", input: "one\ntwo\nthree\nfour\nfive\n", wants: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			got := counter.GetCounts(r).Lines
			if got != tc.wants {
				t.Errorf("Got %d, wants %d", got, tc.wants)
			}
		})
	}
}

func TestCountBytes(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wants int
	}{
		{name: "empty string", input: "", wants: 0},
		{name: "single character string", input: "a", wants: 1},
		{name: "one word", input: "hello", wants: 5},
		{name: "all spaces", input: "       ", wants: 7},
		{name: "five words", input: "one two three four five", wants: 23},
		{name: "newlines and words", input: "one\ntwo\nthree\n", wants: 14},
		{name: "newlines, tab and words", input: "one\ntwo\nthree\nfour\tfive\n", wants: 24},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.input)
			got := counter.GetCounts(r).Bytes
			if got != tc.wants {
				t.Errorf("Got %d, wants %d", got, tc.wants)
			}
		})
	}
}

func TestGetCounts(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		wants counter.Counts
	}{
		{name: "five words", input: "one two three four five\n", wants: counter.Counts{
			Lines: 1,
			Words: 5,
			Bytes: 24,
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			got := counter.GetCounts(reader)
			if got != tc.wants {
				t.Errorf("got: %v, wants: %v", got, tc.wants)
			}
		})
	}
}

func TestPrintCounts(t *testing.T) {
	type inputs struct {
		counts   counter.Counts
		filename []string
		options  counter.DisplayOptions
	}

	testCases := []struct {
		name  string
		input inputs
		wants string
	}{
		{
			name: "five words",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: []string{"words.txt"},
				options: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: true,
				},
			},
			wants: "1 5 24 words.txt\n",
		},
		{
			name: "no filename",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 4,
					Bytes: 18,
				},
				options: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: true,
				},
			},
			wants: "1 4 18\n",
		},
		{
			name: "five words show lines",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: []string{"words.txt"},
				options: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: false,
					ShowBytes: false,
				},
			},
			wants: "1 words.txt\n",
		},
		{
			name: "five words show words",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: []string{"words.txt"},
				options: counter.DisplayOptions{
					ShowLines: false,
					ShowWords: true,
					ShowBytes: false,
				},
			},
			wants: "5 words.txt\n",
		},
		{
			name: "five words show bytes",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: []string{"words.txt"},
				options: counter.DisplayOptions{
					ShowLines: false,
					ShowWords: false,
					ShowBytes: true,
				},
			},
			wants: "24 words.txt\n",
		},
		{
			name: "five words show lines and bytes",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: []string{"words.txt"},
				options: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: false,
					ShowBytes: true,
				},
			},
			wants: "1 24 words.txt\n",
		},
		{
			name: "five words show words and bytes",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: []string{"words.txt"},
				options: counter.DisplayOptions{
					ShowLines: false,
					ShowWords: true,
					ShowBytes: true,
				},
			},
			wants: "5 24 words.txt\n",
		},
		{
			name: "five words show lines and words",
			input: inputs{
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
				},
				filename: []string{"words.txt"},
				options: counter.DisplayOptions{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: false,
				},
			},
			wants: "1 5 words.txt\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			tc.input.counts.Print(&buffer, tc.input.options, tc.input.filename...)
			got := buffer.String()
			if got != tc.wants {
				t.Errorf("got: %v, wants: %v", buffer.Bytes(), []byte(tc.wants))
			}
		})
	}

}

func TestAddCounts(t *testing.T) {
	testCases := []struct {
		name  string
		input []counter.Counts
		wants counter.Counts
	}{
		{
			name: "0 count",
			input: []counter.Counts{
				{
					Lines: 0,
					Words: 0,
					Bytes: 0,
				},
				{
					Lines: 0,
					Words: 0,
					Bytes: 0,
				},
			},
			wants: counter.Counts{
				Lines: 0,
				Words: 0,
				Bytes: 0,
			},
		},
		{
			name: "simple count",
			input: []counter.Counts{
				{
					Lines: 1,
					Words: 2,
					Bytes: 3,
				},
			},
			wants: counter.Counts{
				Lines: 1,
				Words: 2,
				Bytes: 3,
			},
		},
		{
			name: "multi count",
			input: []counter.Counts{
				{
					Lines: 1,
					Words: 2,
					Bytes: 3,
				},
				{
					Lines: 4,
					Words: 5,
					Bytes: 6,
				},
			},
			wants: counter.Counts{
				Lines: 5,
				Words: 7,
				Bytes: 9,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totals := counter.Counts{}
			for _, input := range tc.input {
				totals = totals.Add(input)
			}
			if totals != tc.wants {
				t.Errorf("got: %v wants: %v", totals, tc.wants)
			}
		})
	}
}
