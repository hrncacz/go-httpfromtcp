package server

import (
	"fmt"
	"log"
	"net"

	"github.com/hrncacz/go-httpfromtcp/internal/response"
)

type Server struct {
	serverState serverState
	listener    net.Listener
}

type serverState int

const (
	serverListening serverState = iota
	serverClose
)

func Serve(port int) (*Server, error) {
	portString := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", portString)
	if err != nil {
		log.Fatal(err)
	}
	server := &Server{
		serverState: serverListening,
		listener:    l,
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

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	err := response.WriteStatusLine(conn, 200)
	if err != nil {
		log.Fatal(err)
	}
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Close() error {
	s.serverState = serverClose
	return nil
}
