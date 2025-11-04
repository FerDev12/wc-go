package e2e

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName = "wc-go-test"
)

func TestMain(m *testing.M) {
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binName, "../../cmd")

	buf := &bytes.Buffer{}
	cmd.Stderr = buf

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to build binary:", err, buf.String())
		os.Exit(1)
	}

	result := m.Run()

	os.Remove(binName)
	os.Exit(result)
}

func getCommand(args ...string) (*exec.Cmd, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(dir, binName)

	return exec.Command(path, args...), nil
}

func createFile(dname string, content string) (*os.File, error) {
	file, err := os.CreateTemp(dname, "wc-go-test-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		return nil, fmt.Errorf("failed towrite content to file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close file: %w", err)
	}

	return file, nil
}
