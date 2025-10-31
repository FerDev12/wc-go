package counter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"bloom.io/github.com/FerDev12/wc-go/display"
)

type Counts struct {
	lines uint
	words uint
	bytes uint
}

func (c Counts) Add(other Counts) Counts {
	c.lines += other.lines
	c.words += other.words
	c.bytes += other.bytes
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
func CountWords(data io.Reader) uint {
	wordCount := uint(0)

	// Create scanner
	scanner := bufio.NewScanner(data)
	// Defines the function that Split will use
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		wordCount++
	}

	return wordCount
}

func CountLines(r io.Reader) uint {
	linesCount := uint(0)

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

func CountBytes(r io.Reader) uint {
	byteCount, _ := io.Copy(io.Discard, r)
	return uint(byteCount)
}

func getCountsConcurrent(r io.Reader) Counts {
	linesReader, linesWriter := io.Pipe()
	wordsReader, wordsWriter := io.Pipe()
	bytesReader, bytesWriter := io.Pipe()

	w := io.MultiWriter(linesWriter, wordsWriter, bytesWriter)

	linesChan := make(chan uint)
	wordsChan := make(chan uint)
	bytesChan := make(chan uint)

	go func() {
		defer close(linesChan)
		linesChan <- CountLines(linesReader)
	}()

	go func() {
		defer close(wordsChan)
		wordsChan <- CountWords(wordsReader)
	}()

	go func() {
		defer close(bytesChan)
		bytesChan <- CountBytes(bytesReader)
	}()

	io.Copy(w, r)

	linesWriter.Close()
	wordsWriter.Close()
	bytesWriter.Close()

	return Counts{
		<-linesChan,
		<-wordsChan,
		<-bytesChan,
	}
}

func getCountsSinglePass(r io.Reader) Counts {
	res := Counts{}

	isInsideWord := false
	reader := bufio.NewReader(r)

	for {
		r, size, err := reader.ReadRune()

		if err != nil {
			break
		}

		res.bytes += uint(size)

		if r == '\n' {
			res.lines++
		}

		isSpace := unicode.IsSpace(r)

		if !isSpace && !isInsideWord {
			res.words++
		}

		isInsideWord = !isSpace
	}

	return res
}

func GetCounts(r io.Reader) Counts {
	return getCountsSinglePass(r)
}

func (c Counts) Print(w io.Writer, opts display.Options, suffixes ...string) {
	stats := []string{}

	if opts.ShouldShowLines() {
		stats = append(stats, strconv.Itoa(int(c.lines)))
	}
	if opts.ShouldShowWords() {
		stats = append(stats, strconv.Itoa(int(c.words)))
	}
	if opts.ShouldShowBytes() {
		stats = append(stats, strconv.Itoa(int(c.bytes)))
	}

	line := strings.Join(stats, "\t") + "\t"
	suffixStr := strings.Join(suffixes, " ")

	fmt.Fprintf(w, "%s", line)

	if suffixStr != "" {
		fmt.Fprintf(w, " %s", suffixStr)
	}

	fmt.Fprintf(w, "\n")
}
