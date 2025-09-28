package response

import (
	"errors"
	"fmt"
	"io"

	"github.com/hrncacz/go-httpfromtcp/internal/headers"
)

var statusCode = map[int]string{
	200: "OK",
	400: "Bad Request",
	500: "Server Error",
}

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

type Writer struct {
	writerState writerState
	writer      io.Writer
}

func (w *Writer) WriteStatusLine(status int) error {
	if w.writerState != writerStateStatusLine {
		return errors.New("status line alredy set")
	}
	if err := WriteStatusLine(w.writer, status); err != nil {
		return err
	}
	w.writerState = writerStateHeaders
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState < writerStateHeaders {
		return errors.New("writer currently expects status line")
	} else if w.writerState > writerStateHeaders {
		return errors.New("headers already set")
	}
	if err := WriteHeaders(w.writer, headers); err != nil {
		return err
	}
	w.writerState = writerStateBody
	return nil
}

func (w *Writer) WriteBody(body []byte) (int, error) {
	if w.writerState < writerStateBody {
		return 0, errors.New("writer expects status line and headers before sending body")
	}
	n, err := w.writer.Write(body)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func WriteStatusLine(w io.Writer, status int) error {
	startLine := "HTTP/1.1 "
	statusText, exist := statusCode[status]
	if !exist {
		startLine = fmt.Sprintf("%s%d", startLine, status)
		w.Write([]byte(startLine))
	}
	startLine = fmt.Sprintf("%s%d %s\r\n", startLine, status, statusText)
	w.Write([]byte(startLine))
	return nil
}

func GetDefaultHeaders(contentLength int) headers.Headers {
	responseHeaders := headers.NewHeaders()
	responseHeaders.Set("Content-Length", fmt.Sprintf("%d", contentLength))
	responseHeaders.Set("Connection", "close")
	responseHeaders.Set("Content-Type", "text/plain")
	return responseHeaders
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		responseString := fmt.Sprintf("%s: %s\r\n", k, v)
		w.Write([]byte(responseString))
	}
	w.Write([]byte("\r\n"))
	return nil
}

func NewResponse(w io.Writer) *Writer {
	return &Writer{
		writerState: writerStateStatusLine,
		writer:      w,
	}
}
