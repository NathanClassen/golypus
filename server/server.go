package server

import (
	"fmt"
	"net"
	"strings"
)

type Server struct {
	Listener	net.Listener
	Config		ServerConfiguration
}

type ServerConfiguration struct {
	Address string
	Port string
	Proto string
}

func NewServer(config ServerConfiguration) *Server {

	loc := strings.Join([]string{config.Address, config.Port}, ":")

	l, _ := net.Listen(config.Proto, loc)	//	TODO: handle error

	return &Server{l, config}
}

func (s *Server) Start() {
	for {
		c, _ := s.acceptConnection()	//	TODO: handle error
		s.handleConnection(c)
	}
}

func (s *Server) handleConnection(c net.Conn) {
	
	request := BuildRequest(c)

	fmt.Printf("request line: %s\nheaders: %v\n", request.RequestLine, request.Headers)

	c.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	c.Close()
}

func (s *Server) acceptConnection() (net.Conn, error) {

	c, err := s.Listener.Accept()
	
	if err != nil {
		return nil, err
	}

	return c, nil
}