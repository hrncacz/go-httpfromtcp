package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filepath, err := filepath.Abs("./message.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch := getLinesChannel(file)
	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		defer f.Close()
		currentLine := ""
		for {
			readBuffer := make([]byte, 8)
			_, err := f.Read(readBuffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
			}
			stringFromBuffer := string(readBuffer)
			stringArr := strings.Split(stringFromBuffer, "\n")
			currentLine = lineSeparator(stringArr, currentLine, ch)
		}
	}()
	return ch
}

func lineSeparator(stringArr []string, currentLine string, ch chan string) string {
	if len(stringArr) == 1 {
		return currentLine + stringArr[0]
	} else {
		ch <- currentLine + stringArr[0]
		return lineSeparator(stringArr[1:], "", ch)
	}
}
