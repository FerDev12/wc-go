package e2e

import (
	"fmt"
	"os"
	"testing"

	"bloom.io/github.com/FerDev12/wc-go/test/assert"
)

func TestMultiFile(t *testing.T) {
	dname, err := os.MkdirTemp("", "multi-file-test")
	if err != nil {
		t.Fatal("failed to create directory:", err)
	}

	defer os.RemoveAll(dname)

	fileA, err := createFile(dname, "one two three four five\n")
	if err != nil {
		t.Fatal("failed to create fileA:", err)
	}

	fileB, err := createFile(dname, "foo bar baz\n\n")
	if err != nil {
		t.Fatal("failed to create fileB:", err)
	}

	fileC, err := createFile(dname, "")
	if err != nil {
		t.Fatal("failed to create fileC:", err)
	}

	cmd, err := getCommand(fileA.Name(), fileB.Name(), fileC.Name())
	if err != nil {
		t.Fatal("failed to create command:", err)
	}

	stdout, err := cmd.Output()
	if err != nil {
		t.Fatal("failed to run command:", err)
	}

	got := string(stdout)
	wants := fmt.Sprintf(`    1    5    24 %s
    2    3    13 %s
    0    0     0 %s
    3    8    37 total
`, fileA.Name(), fileB.Name(), fileC.Name())

	assert.Equal(t, wants, got)
}
