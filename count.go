package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

func (c Counts) Print(w io.Writer, opts DisplayOptions, suffixes ...string) {
	stats := []string{}
	showAll := opts.ShouldShowAll()

	if opts.ShowLines || showAll {
		stats = append(stats, strconv.Itoa(c.Lines))
	}
	if opts.ShowWords || showAll {
		stats = append(stats, strconv.Itoa(c.Words))
	}
	if opts.ShowBytes || showAll {
		stats = append(stats, strconv.Itoa(c.Bytes))
	}

	line := strings.Join(stats, "\t") + "\t"
	suffixStr := strings.Join(suffixes, " ")

	fmt.Fprintf(w, "%s", line)

	if suffixStr != "" {
		fmt.Fprintf(w, " %s", suffixStr)
	}

	fmt.Fprintf(w, "\n")
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

// --------- PREVIOUS IMPL ----------
// func CountWords(data []byte) int {
// 	return len(bytes.Fields(data))
// }

// func CountWordsInFile(file *os.File) int {
// 	const BUFFER_SIZE = 2
// 	buffer := make([]byte, BUFFER_SIZE)
// 	leftover := []byte{}

// 	wordCount := 0
// 	isInsideWord := false

// 	for {
// 		size, err := file.Read(buffer)

// 		if err != nil {
// 			if err.Error() == "EOF" {
// 				break
// 			} else {
// 				log.Fatalln("Error reading file: ", err)
// 			}
// 		}

// 		subBuffer := append(leftover, buffer[:size]...)

// 		for len(subBuffer) > 0 {
// 			r, rsize := utf8.DecodeRune(subBuffer)
// 			if r == utf8.RuneError {
// 				break
// 			}

// 			subBuffer = subBuffer[rsize:]

// 			if !unicode.IsSpace(r) && !isInsideWord {
// 				wordCount++
// 			}

// 			isInsideWord = !unicode.IsSpace(r)
// 		}

// 		leftover = leftover[:0]
// 		leftover = append(leftover, subBuffer...)
// 	}

// 	fmt.Print("\n")

// 	return wordCount
// }
