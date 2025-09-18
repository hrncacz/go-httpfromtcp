package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Println(line)
		}
	}

}

func getLinesChannel(c net.Conn) <-chan string {
	ch := make(chan string)
	go func() {
		defer c.Close()
		defer close(ch)
		currentLine := ""
		for {
			readBuffer := make([]byte, 8)
			_, err := c.Read(readBuffer)
			if err != nil {
				if currentLine != "" {
					ch <- currentLine
				}
				if errors.Is(err, io.EOF) {
					break
				}
				break
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
