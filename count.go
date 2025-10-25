package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Counts struct {
	Lines int
	Words int
	Bytes int
}

func (c Counts) String() string {
	return fmt.Sprintf("%d %d %d", c.Lines, c.Words, c.Bytes)
}

func (c Counts) Print(w io.Writer, filenames ...string) {
	fmt.Fprintf(w, "%s", c)

	for _, filename := range filenames {
		fmt.Fprintf(w, " %s", filename)
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

func GetCounts(file io.ReadSeeker) Counts {
	const OFFSET_START = 0

	lines := CountLines(file)
	file.Seek(OFFSET_START, io.SeekStart)

	words := CountWords(file)
	file.Seek(OFFSET_START, io.SeekStart)

	bytes := CountBytes(file)

	return Counts{
		Lines: lines,
		Words: words,
		Bytes: bytes,
	}
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
