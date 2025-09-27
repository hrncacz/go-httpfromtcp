package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hrncacz/go-httpfromtcp/internal/request"
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
		// ch := getLinesChannel(conn)
		// for line := range ch {
		// 	fmt.Println(line)
		// }

		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf(`Request line:
- Method: %s
- Target: %s
- Version: %s
Headers:
`, req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		for k, v := range req.Headers {
			fmt.Printf("- %s: %s\n", k, v)
		}
		body := req.Body
		strBody := string(body)
		fmt.Printf("Body:\n%s", strBody)
	}

}

// func getLinesChannel(c net.Conn) <-chan string {
// 	ch := make(chan string)
// 	go func() {
// 		defer c.Close()
// 		defer close(ch)
// 		currentLine := ""
// 		for {
// 			readBuffer := make([]byte, 8)
// 			_, err := c.Read(readBuffer)
// 			if err != nil {
// 				if currentLine != "" {
// 					ch <- currentLine
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				break
// 			}
// 			stringFromBuffer := string(readBuffer)
// 			stringArr := strings.Split(stringFromBuffer, "\n")
// 			currentLine = lineSeparator(stringArr, currentLine, ch)
// 		}
// 	}()
// 	return ch
// }
//
// func lineSeparator(stringArr []string, currentLine string, ch chan string) string {
// 	if len(stringArr) == 1 {
// 		return currentLine + stringArr[0]
// 	} else {
// 		ch <- currentLine + stringArr[0]
// 		return lineSeparator(stringArr[1:], "", ch)
// 	}
// }
