package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
	if req.RequestLine.RequestTarget == "/yourproblem" {
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
	} else if req.RequestLine.RequestTarget == "/myproblem" {
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
	} else {
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
