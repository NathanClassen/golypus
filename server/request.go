package server

import (
	"bufio"
	"fmt"
	"net"
	"nhclassen/golypus/utils"
	"slices"
	"strings"
)

type Request struct {
	RequestLine string
	Method string
	Resource string
	Headers map[string]string
	Proto string
}

var SUPPORTED_METHODS = []string{
	"GET",
}

func BuildRequest(c net.Conn) Request {
	request := Request{}

	crdr := bufio.NewReader(c)

	data := make([]byte,crdr.Size())

	crdr.Read(data)

	messageLines, _ := utils.MessageLines(data)	//	TODO: handle error
	
	requestLine := messageLines[0]	

	request.RequestLine = requestLine

	headerLines := messageLines[1:]

	headerMap, _ := utils.Headers(headerLines)

	request.Headers = headerMap

	request.setMethod()
	request.setResource()
	request.setProto()

	return request
}

func (r *Request) setResource() {
	parts := strings.Split(r.RequestLine, " ")
	r.Resource = parts[1]


}

func (r *Request) setMethod() error {
	parts := strings.Split(r.RequestLine, " ")
	method := parts[0]

	if !slices.Contains[[]string, string](SUPPORTED_METHODS, method) {
		return fmt.Errorf("method '%s' not supported by Golypus\n", method)
	}

	r.Method = method

	return nil
}

func (r *Request) setProto() error {
	parts := strings.Split(r.RequestLine, " ")
	proto := parts[2]

	if proto != "HTTP/1.1" {
		return fmt.Errorf("Golypus supports HTTP/1.1; protocol %s not supported", proto)
	}

	r.Proto = proto

	return nil
}