package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

const bufferSize = 8

type ParserState int

const (
	ParserInitialized ParserState = iota
	ParserDone
)

type Request struct {
	RequestLine RequestLine
	State       ParserState
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case ParserInitialized:
		n, requestLine, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}

		r.RequestLine = requestLine
		r.State = ParserDone
		return n, nil
	case ParserDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}

}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{State: ParserInitialized}

	for req.State != ParserDone {
		// Resize buffer to twice current size if full
		if readToIndex >= cap(buf) {
			new_buf := make([]byte, len(buf)*2)
			copy(new_buf, buf)
			buf = new_buf
		}

		// Read until the buffer is filled starting at readIndex
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = ParserDone
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		// Try to parse the request
		numBytesParsed, err := req.parse(buf)
		if err != nil {
			return nil, err
		}
		// Copy the buffer if we parsed bytes
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return req, nil
}

func isCapitalized(s string) bool {
	return s == strings.ToUpper(s)
}

func parseRequestLine(text string) (int, RequestLine, error) {
	if !strings.Contains(text, "\r\n") {
		return 0, RequestLine{}, nil
	}

	line := strings.Split(text, "\r\n")[0]
	n := len(line)
	tokens := strings.Split(string(line), " ")

	// Get the Method
	method := tokens[0]
	if !isCapitalized(method) {
		return 0, RequestLine{}, fmt.Errorf("method is not capitalized")
	}

	// Get the Target
	target := tokens[1]

	// Get the version
	// Remove HTTP/ from the version
	version := strings.ReplaceAll(tokens[2], "HTTP/", "")
	if !strings.Contains(version, "1.1") {
		return 0, RequestLine{}, fmt.Errorf("invalid http version: %s", version)
	}

	return n, RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil
}

func PrintRequestLine(rl RequestLine) {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", rl.Method)
	fmt.Printf("- Target: %s\n", rl.RequestTarget)
	fmt.Printf("- Version: %s\n", rl.HttpVersion)
}
