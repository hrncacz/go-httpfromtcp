package response

import (
	"fmt"
	"io"

	"github.com/hrncacz/go-httpfromtcp/internal/headers"
)

var statusCode = map[int]string{
	200: "OK",
	400: "Bad Request",
	500: "Server Error",
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
