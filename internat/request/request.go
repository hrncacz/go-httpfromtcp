package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error io.ReadAll: %s", err)
		return &Request{}, err
	}
	dataString := string(data)
	dataLines := strings.Split(dataString, "\r\n")
	if len(dataLines) == 0 {
		return &Request{}, errors.New("http request string not found")
	}
	requestLineArray := strings.Split(dataLines[0], " ")
	if len(requestLineArray) != 3 {
		return &Request{}, errors.New("not valid request line")
	}
	var httpVersion string
	if requestLineArray[2] == "HTTP/1.1" {
		httpVersion = "1.1"
	} else {
		return &Request{}, errors.New("not supported HTTP version")
	}
	requestLine := RequestLine{
		Method:        requestLineArray[0],
		RequestTarget: requestLineArray[1],
		HttpVersion:   httpVersion,
	}
	if requestLine.HttpVersion != "HTTP/1.1" {
	}
	request := &Request{
		RequestLine: requestLine,
	}
	return request, nil
}
