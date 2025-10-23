package main_test

import (
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
		wants counter.CountsResult
	}{
		{name: "five words", input: "one two three four five\n", wants: counter.CountsResult{
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
				t.Errorf("got %v, wants %v", got, tc.wants)
			}
		})
	}
}
