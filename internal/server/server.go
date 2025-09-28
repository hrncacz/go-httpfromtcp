package server

import (
	"fmt"
	"log"
	"net"

	"github.com/hrncacz/go-httpfromtcp/internal/request"
	"github.com/hrncacz/go-httpfromtcp/internal/response"
)

type Server struct {
	serverState serverState
	listener    net.Listener
	handler     Handler
}

type serverState int

const (
	serverListening serverState = iota
	serverClose
)

type Handler func(w *response.Writer, req *request.Request) *HandlerError

type HandlerError struct {
	StatusCode  int
	ContentType string
	Message     string
}

func Serve(port int, handlerFunction Handler) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatal(err)
	}
	server := &Server{
		serverState: serverListening,
		listener:    l,
		handler:     handlerFunction,
	}
	server.listen()
	return server, nil

}

func (s *Server) listen() {
	defer s.listener.Close()
	fmt.Println("Listening...")
	for s.serverState != serverClose {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.handle(conn)
	}

}

func (s *Server) writeError(conn net.Conn, handleError *HandlerError) error {
	defer conn.Close()
	contentLength := len(handleError.Message)
	headers := response.GetDefaultHeaders(contentLength)
	if handleError.ContentType != "" {
		headers.SetNew("Content-type", handleError.ContentType)
	}
	if err := response.WriteStatusLine(conn, handleError.StatusCode); err != nil {
		return err
	}
	if err := response.WriteHeaders(conn, headers); err != nil {
		return err
	}
	conn.Write([]byte(handleError.Message))
	return nil
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
	}
	buf := response.NewResponse(conn)
	if handlerError := s.handler(buf, req); handlerError != nil {
		s.writeError(conn, handlerError)
	}
}

func (s *Server) Close() error {
	s.serverState = serverClose
	return nil
}
