package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to read from reader: %w", err)
	}

	// Split into parts delimited by \r\n crlf
	parts := strings.Split(string(data), "\r\n")

	requestLine, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, fmt.Errorf("unable to parse request line: %w", err)
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func isCapitalized(s string) bool {
	return s == strings.ToUpper(s)
}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")

	// Get the Method
	method := parts[0]
	if !isCapitalized(method) {
		return RequestLine{}, fmt.Errorf("method is not capitalized")
	}

	// Get the Target
	target := parts[1]

	// Get the version
	// Remove HTTP/ from the version
	version := strings.ReplaceAll(parts[2], "HTTP/", "")
	if !strings.Contains(version, "1.1") {
		return RequestLine{}, fmt.Errorf("invalid http version: %s", version)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   version,
	}, nil
}
