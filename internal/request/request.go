package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/hrncacz/go-httpfromtcp/internal/headers"
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
	Headers     headers.Headers
	Body        []byte
}

type parseState int

const (
	parseStateInitialized parseState = iota
	parseStateHeaders
	parseStateBody
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
		r.ParseState = parseStateHeaders
		return n, nil
	case parseStateHeaders:
		fallthrough
	case parseStateBody:
		totalBytesParsed := 0
		for r.ParseState != parseStateDone {
			n, err := r.parseSingle(data[totalBytesParsed:])
			if err != nil {
				return totalBytesParsed, err
			}
			totalBytesParsed += n
			if n == 0 {
				break
			}
		}
		return totalBytesParsed, nil
	case parseStateDone:
		return 0, fmt.Errorf("error reading data from done state")
	default:

		return 0, fmt.Errorf("error reading data from done state")
	}
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.ParseState {
	case parseStateHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.ParseState = parseStateBody
			return n, nil
		} else {
			return n, nil
		}

	case parseStateBody:
		value, exist := r.Headers.Get("Content-Length")
		if !exist {
			r.Body = slices.Concat(r.Body, data)
			r.ParseState = parseStateDone
			return 0, nil
		}
		contentLength, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid content-length header: %s", err)
		}
		r.Body = slices.Concat(r.Body, data)
		if len(r.Body) > contentLength {
			return len(data), errors.New("length of body is greater thatn value in content-length header")
		} else if len(r.Body) == contentLength {
			r.ParseState = parseStateDone
			return len(data), nil
		} else {
			return len(data), nil
		}
	case parseStateDone:
		return 0, fmt.Errorf("error reading data from done state")
	default:

		return 0, fmt.Errorf("error reading data from done state")
	}
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := &Request{
		ParseState: parseStateInitialized,
		Body:       []byte{},
	}
	readToIndex := 0
	request.Headers = headers.NewHeaders()
	buf := make([]byte, buffSize)
	for request.ParseState != parseStateDone {
		if readToIndex >= len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}
		nrBytesToRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) && request.ParseState == parseStateDone {
				return request, nil
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
