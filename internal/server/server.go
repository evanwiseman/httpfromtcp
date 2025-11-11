package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/evanwiseman/httpfromtcp/internal/request"
	"github.com/evanwiseman/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	isOpen   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listener,
		handler:  handler,
	}
	server.isOpen.Store(true)
	go server.listen()

	return server, nil
}

func (s *Server) Close() {
	s.isOpen.Store(false)
	s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.isOpen.Load() {
				return
			}
			log.Println("error accepting connection", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	w := response.NewWriter(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}

	s.handler(w, req)
}

type Handler func(w *response.Writer, req *request.Request)
