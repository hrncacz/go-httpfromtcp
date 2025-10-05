package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"syscall"

	"github.com/hrncacz/go-httpfromtcp/internal/headers"
	"github.com/hrncacz/go-httpfromtcp/internal/request"
	"github.com/hrncacz/go-httpfromtcp/internal/response"
	"github.com/hrncacz/go-httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handlerFunc)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handlerFunc(w *response.Writer, req *request.Request) *server.HandlerError {
	switch req.RequestLine.RequestTarget {
	case "yourproblem":
		errorMessage := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`
		return &server.HandlerError{
			StatusCode:  400,
			Message:     errorMessage,
			ContentType: "text/html",
		}
	case "myproblem":
		errorMessage := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`
		return &server.HandlerError{
			StatusCode:  500,
			Message:     errorMessage,
			ContentType: "text/html",
		}
	case "/httpbin/html":
		header := response.GetDefaultHeaders(0)
		header.Remove("Content-Length")
		header.SetNew("Transfer-Encoding", "chunked")
		header.SetNew("Trailer", "X-Content-SHA256, X-Content-Length")
		res, err := http.Get("https://httpbin.org/html")
		if err != nil {
			return &server.HandlerError{
				StatusCode:  500,
				Message:     "Unable to get data from remote server",
				ContentType: "text/html",
			}

		}
		w.WriteStatusLine(200)
		w.WriteHeaders(header)
		defer res.Body.Close()
		acceptedBytes := []byte{}
		acceptedLength := 0
		for {
			buf := make([]byte, 1024)
			readBytes, err := res.Body.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					_, innerError := w.WriteChunkedBodyDone(true)
					if innerError != nil {
						return &server.HandlerError{
							StatusCode:  500,
							Message:     "Unable to close chunked body",
							ContentType: "text/html",
						}
					}
					sha256Data := sha256.Sum256(acceptedBytes)
					hashString := hex.EncodeToString(sha256Data[:])
					trailers := headers.NewHeaders()
					trailers.SetNew("X-Content-Length", strconv.Itoa(acceptedLength))
					trailers.SetNew("X-Content-SHA256", hashString)
					w.WriteTrailers(trailers)
					return nil
				}
				return &server.HandlerError{
					StatusCode:  500,
					Message:     "Issue with reading data from remote server",
					ContentType: "text/html",
				}
			}
			_, err = w.WriteChunkedBody(buf[:readBytes])
			if err != nil {
				return &server.HandlerError{
					StatusCode:  500,
					Message:     "Issue with writing body",
					ContentType: "text/html",
				}
			}
			acceptedBytes = slices.Concat(acceptedBytes, buf[:readBytes])
			acceptedLength += readBytes
		}
	case "/video":
		videoPath, err := filepath.Abs("/home/martin/bootDev/boot_dev_httpfromtcp/assets/vim.mp4")
		if err != nil {
			return &server.HandlerError{
				StatusCode:  500,
				Message:     "Unable to get data from remote server",
				ContentType: "text/html",
			}

		}
		fmt.Println(videoPath)
		data, err := os.ReadFile(videoPath)
		if err != nil {
			return &server.HandlerError{
				StatusCode:  500,
				Message:     "Unable to get data from remote server",
				ContentType: "text/html",
			}

		}
		header := response.GetDefaultHeaders(len(data))
		header.SetNew("Content-type", "video/mp4")
		w.WriteStatusLine(200)
		w.WriteHeaders(header)
		_, err = w.WriteBody(data)
		if err != nil {
			return &server.HandlerError{
				StatusCode:  500,
				Message:     "Issue with writing body",
				ContentType: "text/html",
			}
		}
		return nil
	default:
		w.WriteStatusLine(200)
		responseMessage := []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`)

		headers := response.GetDefaultHeaders(len(responseMessage))
		headers.SetNew("Content-type", "text/html")
		w.WriteHeaders(headers)
		w.WriteBody(responseMessage)

		return nil
	}
}
