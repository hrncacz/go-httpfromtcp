package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strings"
)

const (
	buffSize = 8
	crlf     = "\r\n"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	ParseState  parseState
	RequestLine RequestLine
}

type parseState int

const (
	parseStateInitialized parseState = iota
	parseStateDone
)

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParseState {
	case parseStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		} else if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.ParseState = parseStateDone
		return n, nil
	case parseStateDone:
		return 0, fmt.Errorf("error reading data from done state")
	default:

		return 0, fmt.Errorf("error reading data from done state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{}
	readToIndex := 0
	buf := make([]byte, buffSize)
	for request.ParseState != parseStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		nrBytesToRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.ParseState = parseStateDone
				break
			}
			return nil, err
		}
		readToIndex += nrBytesToRead
		numBytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineString := string(data[:idx])
	requestLine, err := parseRequestLineString(requestLineString)
	if err != nil {
		return nil, 0, err
	}
	return requestLine, idx + 2, nil
}

func parseRequestLineString(str string) (*RequestLine, error) {
	requestLineArray := strings.Split(str, " ")
	if len(requestLineArray) != 3 {
		return nil, errors.New("not valid request line")
	}
	if !isValidMethod(requestLineArray[0]) {
		return nil, errors.New("invalid http method")
	}
	_, err := url.ParseRequestURI(requestLineArray[1])
	if err != nil {
		return nil, errors.New("invalid url path")
	}
	var httpVersion string
	if requestLineArray[2] == "HTTP/1.1" {
		httpVersion = "1.1"
	} else {
		return nil, errors.New("not supported HTTP version")
	}
	requestLine := &RequestLine{
		Method:        requestLineArray[0],
		RequestTarget: requestLineArray[1],
		HttpVersion:   httpVersion,
	}

	return requestLine, nil
}

func isValidMethod(method string) bool {
	httpMethods := []string{"GET", "POST", "PUT", "DELETE"}
	return slices.Contains(httpMethods, method)
}
