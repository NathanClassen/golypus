package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

/*
	Headers parses header fieldnames and fieldvalues from each
	headerLines and creates a map of headers.

	Per RFC9112, whitespace between any header fieldname and the
	proceeding colon is not allowed and will cause err !=nil
	https://www.rfc-editor.org/rfc/rfc9112#name-field-syntax
*/
func Headers(headerLines []string) (map[string]string, error) {
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
	MessageLines parses out the request line and header field
	lines into a slice of strings.

	Per RFC9112, an empty line (CRLF - \r\n) before the request
	line will be ignored.
	https://www.rfc-editor.org/rfc/rfc9112#name-message-parsing
*/
func MessageLines(d []byte) ([]string, error) {
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

	//	allow for empty line at beginning of message, before the request line
	if d[0] == '\r' && d[1] == '\n' {
		d = d[2:]
	}

	for {
		if d[0] == '\r' && d[1] == '\n' {
			break
		}
		line, startOfNextLine, _ := SingleLine(d)
		lines = append(lines, line)
		d = d[startOfNextLine:]
	}
	
	return lines, nil
}

func SingleLine(d []byte) (string, int, error) {
	lineObtained := false
	line := []byte{}
	var startOfNextLine int

	for i, b := range(d) {
		if b == '\r' {
			if (d[i+1] != '\n') {
				return "", i, errors.New("bare cr encountered")
			} else {
				lineObtained = true
				startOfNextLine = i+2
			}
		} else {
			line = append(line, b)
		}

		if lineObtained {
			break
		}
	}

	return string(line), startOfNextLine, nil
}