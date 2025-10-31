package counter

import (
	"bytes"
	"strings"
	"testing"

	"bloom.io/github.com/FerDev12/wc-go/display"
)

func TestCountWords(t *testing.T) {
	var testCases = []struct {
		name  string
		input string
		wants uint
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
			got := GetCounts(reader).words
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
		wants uint
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
			got := GetCounts(r).lines
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
		wants uint
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
			got := GetCounts(r).bytes
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
		wants Counts
	}{
		{name: "five words", input: "one two three four five\n", wants: Counts{
			lines: 1,
			words: 5,
			bytes: 24,
		}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			got := GetCounts(reader)
			if got != tc.wants {
				t.Errorf("got: %v, wants: %v", got, tc.wants)
			}
		})
	}
}

func TestAddCounts(t *testing.T) {
	testCases := []struct {
		name  string
		input []Counts
		wants Counts
	}{
		{
			name: "0 count",
			input: []Counts{
				{
					lines: 0,
					words: 0,
					bytes: 0,
				},
				{
					lines: 0,
					words: 0,
					bytes: 0,
				},
			},
			wants: Counts{
				lines: 0,
				words: 0,
				bytes: 0,
			},
		},
		{
			name: "simple count",
			input: []Counts{
				{
					lines: 1,
					words: 2,
					bytes: 3,
				},
			},
			wants: Counts{
				lines: 1,
				words: 2,
				bytes: 3,
			},
		},
		{
			name: "multi count",
			input: []Counts{
				{
					lines: 1,
					words: 2,
					bytes: 3,
				},
				{
					lines: 4,
					words: 5,
					bytes: 6,
				},
			},
			wants: Counts{
				lines: 5,
				words: 7,
				bytes: 9,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			totals := Counts{}
			for _, input := range tc.input {
				totals = totals.Add(input)
			}
			if totals != tc.wants {
				t.Errorf("got: %v wants: %v", totals, tc.wants)
			}
		})
	}
}

func TestPrintCounts(t *testing.T) {
	type inputs struct {
		counts   Counts
		filename []string
		options  display.Options
	}

	testCases := []struct {
		name  string
		input inputs
		wants string
	}{
		{
			name: "five words",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 5,
					bytes: 24,
				},
				filename: []string{"words.txt"},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: true,
				}),
			},
			wants: "1\t5\t24\t words.txt\n",
		},
		{
			name: "no filename",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 4,
					bytes: 18,
				},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: true,
				}),
			},
			wants: "1\t4\t18\t\n",
		},
		{
			name: "five words show lines",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 5,
					bytes: 24,
				},
				filename: []string{"words.txt"},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: true,
					ShowWords: false,
					ShowBytes: false,
				}),
			},
			wants: "1\t words.txt\n",
		},
		{
			name: "five words show words",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 5,
					bytes: 24,
				},
				filename: []string{"words.txt"},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: false,
					ShowWords: true,
					ShowBytes: false,
				}),
			},
			wants: "5\t words.txt\n",
		},
		{
			name: "five words show bytes",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 5,
					bytes: 24,
				},
				filename: []string{"words.txt"},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: false,
					ShowWords: false,
					ShowBytes: true,
				}),
			},
			wants: "24\t words.txt\n",
		},
		{
			name: "five words show lines and bytes",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 5,
					bytes: 24,
				},
				filename: []string{"words.txt"},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: true,
					ShowWords: false,
					ShowBytes: true,
				}),
			},
			wants: "1\t24\t words.txt\n",
		},
		{
			name: "five words show words and bytes",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 5,
					bytes: 24,
				},
				filename: []string{"words.txt"},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: false,
					ShowWords: true,
					ShowBytes: true,
				}),
			},
			wants: "5\t24\t words.txt\n",
		},
		{
			name: "five words show lines and words",
			input: inputs{
				counts: Counts{
					lines: 1,
					words: 5,
					bytes: 24,
				},
				filename: []string{"words.txt"},
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines: true,
					ShowWords: true,
					ShowBytes: false,
				}),
			},
			wants: "1\t5\t words.txt\n",
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

var benchData = []string{
	"this is a test data string that runs across multiple lines.",
	"one two three four five six seven eight",
	"this is a weird string.",
}

func BenchmarkGetCounts(b *testing.B) {
	i := 0
	for b.Loop() {
		data := benchData[i%len(benchData)]
		r := strings.NewReader(data)
		getCountsConcurrent(r)
		i++
	}
}

func BenchmarkGetCountsSinglePass(b *testing.B) {
	i := 0
	for b.Loop() {
		data := benchData[i%len(benchData)]
		r := strings.NewReader(data)
		getCountsSinglePass(r)
		i++
	}
}
