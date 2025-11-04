package e2e

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
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

	if got != wants {
		t.Errorf("stdout is not correct: got: %s, wants: %s", got, wants)
	}
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

	if got != wants {
		t.Errorf("stdout is not correct: got: %s, wants: %s", got, wants)
	}
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

	if gotStderr != wantsStderr {
		t.Errorf("stderr is not correct:\ngot: %s\nwants: %s", gotStderr, wantsStderr)
	}

	if gotStdout != wantsStdout {
		t.Errorf("stdout is not correct:\ngot: %s\nwants: %s", gotStdout, wantsStdout)
	}
}
