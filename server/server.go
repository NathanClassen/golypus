package server

import (
	"bufio"
	"fmt"
	"net"
	"nhclassen/golypus/utils"
	"strings"
)

/*
	create a server object with a listener
	- should have a start method to begin accepting connections
	- should have some method about handling connections
	- should accept connections in concurrent fashion
*/

type Server struct {
	Listener	net.Listener
	Config		ServerConfiguration
}

type ServerConfiguration struct {
	Address string
	Port string
	Proto string
}

type Request struct {
	RequestLine []byte
	Method string
	Resource string
	Headers map[string]string
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
	// rq := Request{}
	/*
		The normal procedure for parsing an HTTP message is to read the start-line into a structure, read
		each header field line into a hash table by field name until the empty line, and then use the
		parsed data to determine if a message body is expected. If a message body has been indicated,
		then it is read as a stream until an amount of octets equal to the message body length is read or
		the connection is closed.

		- read start line into a structure (rq)
		- read headers into a structure
		- determine if a body should be read
		- read body as a stream according to headers


		- get read all bytes by reading lines separated by crlf
			will start by reading lines separated by \n because im not
			sure of the best way to read until \r\n yet
		- put lines into structure
		- parse lines and store lines in string (for starting line) and
			map (for headers)
	*/
	crdr := bufio.NewReader(c)

	data := make([]byte,crdr.Size())

	crdr.Read(data)
	fmt.Println(string(data))
	requestLine, startOfNextLine, _ := utils.RequestLine(data)	//	TODO: handle error

	fmt.Printf("requested: %v\n", requestLine)

	headerLines, _ := utils.HeaderFieldLines(data[startOfNextLine:])	//	TODO: handle error

	fmt.Printf("headers: %v", headerLines)

	headerMap, _ := utils.ParseHeaders(headerLines)
	
	fmt.Printf("header map: %v\n", headerMap)

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