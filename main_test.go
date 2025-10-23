package main_test

import (
	"strings"
	"testing"

	counter "bloom.io/word-counter"
)

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

func TesetCountWords(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reader := strings.NewReader(tc.input)
			got := counter.CountWords(reader)
			if got != tc.wants {
				t.Errorf("got %d, wants %d", got, tc.wants)
			}
		})
	}
}
