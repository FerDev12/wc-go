package display_test

import (
	"bytes"
	"testing"

	"bloom.io/github.com/FerDev12/wc-go/display"
	"bloom.io/github.com/FerDev12/wc-go/test/assert"
)

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
			assert.Equal(t, tc.wants, buffer.String())
		})
	}
}
