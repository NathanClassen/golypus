package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var logger log.Logger

func init() {
	logger = *log.New(os.Stdout, "golypus webserver v1: ", 0)
}

type Server struct {
	Listener net.Listener
	Config   ServerConfiguration
}

type ServerConfiguration struct {
	Address string
	Port    string
	Proto   string
	Static  string
}

func NewServer(config ServerConfiguration) *Server {

	loc := strings.Join([]string{config.Address, config.Port}, ":")

	l, err := net.Listen(config.Proto, loc) //	TODO: handle error
	if err != nil {
		logger.Fatalln("failed to create new server: ", err)
	}

	return &Server{l, config}
}

func (s *Server) Start() {
	fmt.Printf(`

   ____________
  /            |
 /     ________|
|     |  ______
|     | |___   |
|      \____|  |
|              |
|              |
 \_____________| olypus, accepting connections at %s
 
 ================================================
 ================================================

`, s.Listener.Addr())

	for {
		c, err := s.acceptConnection() //	TODO: handle error
		if err != nil {
			logger.Fatalln("error accepting connections: ", err)
		}
		go s.handleConnection(c)
	}
}

func (s *Server) handleConnection(c net.Conn) {

	request := BuildRequest(c)

	logger.Printf("got req : %+v\n", request)

	c.Write([]byte(fmt.Sprintf("%s%s\n", OK200, request.RequestLine)))
	c.Close()
}

func (s *Server) acceptConnection() (net.Conn, error) {

	c, err := s.Listener.Accept()

	if err != nil {
		return nil, err
	}

	return c, nil
}
