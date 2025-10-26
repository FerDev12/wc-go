package counter

import (
	"bufio"
	"io"
	"os"
	"unicode"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

func (c Counts) Add(other Counts) Counts {
	c.Lines += other.Lines
	c.Words += other.Words
	c.Bytes += other.Bytes
	return c
}

func CountFile(filename string) (Counts, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return Counts{}, err
	}

	return GetCounts(file), nil
}

// By making our argument accept any value that conforms to the io.Reader interface
// we are able to accept various data types such as files or a slice of bytes
func CountWords(data io.Reader) int {
	wordCount := 0

	// Create scanner
	scanner := bufio.NewScanner(data)
	// Defines the function that Split will use
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordCount++
	}

	return wordCount
}

func CountLines(r io.Reader) int {
	linesCount := 0

	reader := bufio.NewReader(r)

	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if r == '\n' {
			linesCount++
		}
	}

	return linesCount
}

func CountBytes(r io.Reader) int {
	byteCount, _ := io.Copy(io.Discard, r)
	return int(byteCount)
}

func GetCounts(r io.Reader) Counts {
	res := Counts{}

	isInsideWord := false
	reader := bufio.NewReader(r)

	for {
		r, size, err := reader.ReadRune()

		if err != nil {
			break
		}

		res.Bytes += size

		if r == '\n' {
			res.Lines++
		}

		isSpace := unicode.IsSpace(r)

		if !isSpace && !isInsideWord {
			res.Words++
		}

		isInsideWord = !isSpace
	}

	return res
}
