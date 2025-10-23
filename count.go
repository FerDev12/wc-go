package main

import (
	"bufio"
	"io"
	"os"
)

func CountWordsForFile(filename string) (int, error) {
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		return 0, err
	}

	return CountWords(file), nil
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
