package display_test

import (
	"bytes"
	"testing"

	counter "bloom.io/github.com/FerDev12/wc-go"
	"bloom.io/github.com/FerDev12/wc-go/display"
)

func TestPrintCounts(t *testing.T) {
	type inputs struct {
		counts   counter.Counts
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
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
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
				counts: counter.Counts{
					Lines: 1,
					Words: 4,
					Bytes: 18,
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
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
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
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
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
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
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
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
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
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
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
				counts: counter.Counts{
					Lines: 1,
					Words: 5,
					Bytes: 24,
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
			tc.input.options.PrintCounts(&buffer, tc.input.counts, tc.input.filename...)
			got := buffer.String()
			if got != tc.wants {
				t.Errorf("got: %v, wants: %v", buffer.Bytes(), []byte(tc.wants))
			}
		})
	}
}

func TestPrintHeader(t *testing.T) {
	type inputs struct {
		options display.Options
	}

	testCases := []struct {
		name  string
		input inputs
		wants string
	}{
		{
			name: "show all with header",
			input: inputs{
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines:  true,
					ShowWords:  true,
					ShowBytes:  true,
					ShowHeader: true,
				}),
			},
			wants: "lines\twords\tcharacters\t\n",
		},
		{
			name: "show lines with header",
			input: inputs{
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines:  true,
					ShowWords:  false,
					ShowBytes:  false,
					ShowHeader: true,
				}),
			},
			wants: "lines\t\n",
		},
		{
			name: "show words with header",
			input: inputs{
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines:  false,
					ShowWords:  true,
					ShowBytes:  false,
					ShowHeader: true,
				}),
			},
			wants: "words\t\n",
		},
		{
			name: "show bytes with header",
			input: inputs{
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines:  false,
					ShowWords:  false,
					ShowBytes:  true,
					ShowHeader: true,
				}),
			},
			wants: "characters\t\n",
		},
		{
			name: "show lines and words with header",
			input: inputs{
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines:  true,
					ShowWords:  true,
					ShowBytes:  false,
					ShowHeader: true,
				}),
			},
			wants: "lines\twords\t\n",
		},
		{
			name: "show lines and bytes with header",
			input: inputs{
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines:  true,
					ShowWords:  false,
					ShowBytes:  true,
					ShowHeader: true,
				}),
			},
			wants: "lines\tcharacters\t\n",
		},
		{
			name: "show words and bytes with header",
			input: inputs{
				options: display.NewOptions(display.NewOptionsArgs{
					ShowLines:  false,
					ShowWords:  true,
					ShowBytes:  true,
					ShowHeader: true,
				}),
			},
			wants: "words\tcharacters\t\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buffer := bytes.Buffer{}
			tc.input.options.PrintHeader(&buffer)
			if buffer.String() != tc.wants {
				t.Errorf("got: %v, wants: %v", buffer.Bytes(), []byte(tc.wants))
			}
		})
	}
}
