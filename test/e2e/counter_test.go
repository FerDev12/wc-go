package e2e

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"bloom.io/github.com/FerDev12/wc-go/test/e2e/assert"
)

func TestStdin(t *testing.T) {
	cmd, err := getCommand()

	if err != nil {
		t.Fatal("couldn't get working directory:", err)
	}

	output := &bytes.Buffer{}

	cmd.Stdin = strings.NewReader("one two three\n")
	cmd.Stdout = output

	if err := cmd.Run(); err != nil {
		t.Fatal("failed to run command:", err)
	}

	wants := "    1    3    14\n"
	got := output.String()
	assert.Equal(t, wants, got, "stdout is not correct")
}

func TestSingleFile(t *testing.T) {
	file, err := os.CreateTemp("", "wc-go-test-*")
	if err != nil {
		t.Fatal("couldn't create temp file:", err)
	}

	defer os.Remove(file.Name())

	_, err = file.WriteString("one two three\nfour five six\nseven eight nine\n")
	if err != nil {
		t.Fatal("failed to write to temp file:", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal("failed to close temp file:", err)
	}

	cmd, err := getCommand(file.Name())
	if err != nil {
		t.Fatal("couldn't get working directory:", err)
	}

	output := &bytes.Buffer{}
	cmd.Stdout = output

	if err = cmd.Run(); err != nil {
		t.Fatal("failed to run command:", err)
	}

	wants := fmt.Sprintf("    3    9    45 %s\n    3    9    45 total\n", file.Name())
	got := output.String()
	assert.Equal(t, wants, got, "stdout is not correct")
}

func TestNonExistingFile(t *testing.T) {
	filename := "non-existent.txt"

	cmd, err := getCommand(filename)
	if err != nil {
		t.Fatal("failed to get command:", err)
	}

	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	cmd.Stderr = stderr
	cmd.Stdout = stdout

	if err = cmd.Run(); err == nil {
		t.Error("command succeeded when it shouldn't")
	}

	if err.Error() != "exit status 1" {
		t.Errorf("unexpected error: got: %s, wants: exit status 1", err)
	}

	wantsStderr := fmt.Sprintf("wc-go: open %s: no such file or directory\n", filename)
	wantsStdout := "    0    0    0 total\n"

	gotStderr := stderr.String()
	gotStdout := stdout.String()

	assert.Equal(t, wantsStderr, gotStderr, "stderr is not correct")
	assert.Equal(t, wantsStdout, gotStdout, "stdout is not correct")
}

func TestFlags(t *testing.T) {
	dname, err := os.MkdirTemp("", "flags")
	if err != nil {
		t.Fatal("failed to create temp directory:", err)
	}
	defer os.RemoveAll(dname)

	file, err := createFile(dname, "one two three\nfour five six\n")
	if err != nil {
		t.Fatal("failed to create temp file:", err)
	}

	type inputs struct {
		flags   []string
		content string
	}

	testCases := []struct {
		name  string
		input inputs
		wants string
	}{
		{
			name: "-l (line) flag",
			input: inputs{
				flags:   []string{"-l"},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    2 %s
    2 total
`, file.Name()),
		},
		{
			name: "-w (word) flag",
			input: inputs{
				flags:   []string{"-w"},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    6 %s
    6 total
`, file.Name()),
		},
		{
			name: "-c (bytes) flag",
			input: inputs{
				flags:   []string{"-c"},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    28 %s
    28 total
`, file.Name()),
		},
		{
			name: "-l -w (line and word) flag",
			input: inputs{
				flags:   []string{"-l", "-w"},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    2    6 %s
    2    6 total
`, file.Name()),
		},
		{
			name: "-l -c (line and byte) flag",
			input: inputs{
				flags:   []string{"-l", "-c"},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    2    28 %s
    2    28 total
`, file.Name()),
		},
		{
			name: "-w -c (word and byte) flag",
			input: inputs{
				flags:   []string{"-w", "-c"},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    6    28 %s
    6    28 total
`, file.Name()),
		},
		{
			name: "all flags",
			input: inputs{
				flags:   []string{"-l", "-w", "-c"},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    2    6    28 %s
    2    6    28 total
`, file.Name()),
		},
		{
			name: "no flags",
			input: inputs{
				flags:   []string{},
				content: "one two three\nfour five six\n",
			},
			wants: fmt.Sprintf(`    2    6    28 %s
    2    6    28 total
`, file.Name()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := append(tc.input.flags, file.Name())
			cmd, err := getCommand(args...)
			if err != nil {
				t.Error("failed to get command:", err)
			}

			stdout := &bytes.Buffer{}
			cmd.Stdout = stdout

			if err := cmd.Run(); err != nil {
				t.Error("failed to run command", err)
			}

			got := stdout.String()
			assert.Equal(t, tc.wants, got, "stdout is not correct")
		})
	}
}
