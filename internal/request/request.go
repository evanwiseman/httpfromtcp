package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/evanwiseman/httpfromtcp/internal/headers"
)

const crlf = "\r\n"
const bufferSize = 8

type ParserState int

const (
	ParserInitialized ParserState = iota
	ParserHeaders
	ParserBody
	ParserDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	State       ParserState
}

func (r *Request) parse(data []byte) (n int, err error) {
	totalBytesParsed := 0
	for r.State != ParserDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil { // something went wrong
			return 0, err
		}
		if n == 0 { // need more data
			break
		}
		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (n int, err error) {
	switch r.State {
	case ParserInitialized:
		n, requestLine, err := parseRequestLine(data)
		if err != nil { // something went wrong
			return 0, err
		}
		if n == 0 { // need more data
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = ParserHeaders
		return n, nil
	case ParserHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil { // something went wrong
			return 0, err
		}
		if done {
			r.State = ParserBody
		}
		return n, nil
	case ParserBody:
		lengthStr, ok := r.Headers.Get("content-length")
		if !ok {
			r.State = ParserDone
			return 0, nil
		}
		r.Body = append(r.Body, data...)

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return 0, fmt.Errorf("error: invalid content-length: %w", err)
		}

		if len(r.Body) > length {
			return 0, fmt.Errorf("error: body is larger than content-length %v != %v", len(r.Body), length)
		}
		if len(r.Body) == length {
			r.State = ParserDone
		}

		return len(data), nil
	case ParserDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("error: unknown state")
	}
}

func parseRequestLine(data []byte) (int, *RequestLine, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil, nil
	}
	requestLineText := string(data[:idx])
	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return 0, nil, err
	}
	return idx + 2, requestLine, nil
}

func requestLineFromString(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, fmt.Errorf("invalid method: %s", method)
		}
	}

	requestTarget := parts[1]

	versionParts := strings.Split(parts[2], "/")
	if len(versionParts) != 2 {
		return nil, fmt.Errorf("malformed start-line: %s", str)
	}

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", httpPart)
	}
	version := versionParts[1]
	if version != "1.1" {
		return nil, fmt.Errorf("unrecognized HTTP-version: %s", version)
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: requestTarget,
		HttpVersion:   versionParts[1],
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	req := &Request{
		Headers: headers.NewHeaders(),
		State:   ParserInitialized,
		Body:    make([]byte, 0),
	}

	for req.State != ParserDone {
		// Resize buffer to twice current size if full
		if readToIndex >= cap(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		// Read until the buffer is filled starting at readIndex
		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != ParserDone {
					return nil, fmt.Errorf("incomplete request: %w", err)
				}
				break
			}

			return nil, err
		}
		readToIndex += numBytesRead

		// Try to parse the request
		numBytesParsed, err := req.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		// Copy the buffer if we parsed bytes
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed
	}

	return req, nil
}

func PrintRequest(req *Request) {
	fmt.Println("Request line:")
	fmt.Printf("- Method: %s\n", req.RequestLine.Method)
	fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
	fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)

	fmt.Println("Headers:")
	for k, v := range req.Headers {
		fmt.Printf("- %s: %s\n", k, v)
	}

	fmt.Println("Body:")
	fmt.Println(string(req.Body))
}
