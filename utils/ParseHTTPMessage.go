package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func ParseHeaders(headerLines []string) (map[string]string, error) {
	headers := map[string]string{}

	for _, line := range(headerLines) {
		parts := strings.Split(line, ":")
		fieldname := parts[0]
		fieldvalue := strings.TrimSpace(parts[1])

		if unicode.IsSpace(rune(fieldname[len(fieldname) - 1])) {
			return nil, fmt.Errorf("illegal whitespace found between header fieldname and colon: %s:%s",fieldname,parts[1])
		}

		headers[fieldname] = fieldvalue
	}

	return headers, nil
}

/*
	RequestLine takes a slice of bytes representing the entirety
	of an HTTP request message, and parses out the request line
	returning it as a string and the byte location of the next
	line of the message (the first header line).

	Per RFC9112, an empty line (CRLF - \r\n) before the request
	line will be ignored.
	https://www.rfc-editor.org/rfc/rfc9112#name-message-parsing
*/
func RequestLine(d []byte) (string, int, error) {
	/*

		In the interest of robustness, a server that is expecting to receive and parse a 
		request-line SHOULD ignore at least one empty line (CRLF) received prior to the request-line.

	*/

	if d[0] == '\r' && d[1] == '\n' {
		d = d[2:]
	}

	return ParseLine(d)
}

/*
	HeaderFieldLines takes a slice of bytes representing the en–
	tirety of an HTTP request message, and parses out the header
	fields as lines without evalutation the fieldnames or field
	values

	TODO: this could just be used to parse out all of the non-content
	lines of the request rather than arbitrarily using it only for 
	headers. As it is I remove the request line before calling this,
	but if the whole request was just passed in it would add the request
	line and headers only as lines in the map and I could then work with that
*/
func HeaderFieldLines(d []byte) ([]string, error) {
	/*
		Messages are parsed using a generic algorithm, independent of the individual field names. 
		The contents within a given field line value are not parsed until a later stage of message 
		interpretation (usually after the message's entire field section has been processed).

		No whitespace is allowed between the field name and colon. In the past, differences in the 
		handling of such whitespace have led to security vulnerabilities in request routing and response 
		handling. A server MUST reject, with a response status code of 400 (Bad Request), any received 
		request message that contains whitespace between a header field name and colon. A proxy MUST 
		remove any such whitespace from a response message before forwarding the message downstream.

		A field line value might be preceded and/or followed by optional whitespace (OWS); a single SP 
		preceding the field line value is preferred for consistent readability by humans. The field line 
		value does not include that leading or trailing whitespace: OWS occurring before the first 
		non-whitespace octet of the field line value, or after the last non-whitespace octet of the field 
		line value, is excluded by parsers when extracting the field line value from a field line.
	*/
	lines := []string{}

	for {
		if d[0] == '\r' && d[1] == '\n' {
			break
		}
		headerLine, startOfNextLine, _ := ParseLine(d)
		lines = append(lines, headerLine)
		d = d[startOfNextLine:]
	}
	
	return lines, nil
}

func ParseLine(d []byte) (string, int, error) {
	lineObtained := false
	line := []byte{}
	var startOfNextLine int

	for i, b := range(d) {
		if b == '\r' {
			fmt.Printf("cr at %v\n", i)
			if (d[i+1] != '\n') {
				fmt.Printf("it was a bare cr\n")
				return "", i, errors.New("bare cr encountered")
			} else {
				fmt.Printf("got line\n")
				lineObtained = true
				startOfNextLine = i+2
			}
		} else {
			fmt.Printf("adding %v\n", string(b))
			line = append(line, b)
			fmt.Printf("line is %v\n", string(line))
		}

		if lineObtained {
			break
		}
	}

	fmt.Printf("going to return %v\n", string(line))
	return string(line), startOfNextLine, nil
}