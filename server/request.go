package server

import (
	"bufio"
	"net"
	"nhclassen/golypus/utils"
)

type Request struct {
	RequestLine string
	Method string
	Resource string
	Headers map[string]string
	Proto string
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

	return request
}