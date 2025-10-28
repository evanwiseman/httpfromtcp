package headers

import (
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	text := string(data)
	crlfIdx := strings.Index(text, crlf)
	if crlfIdx == -1 {
		return 0, false, nil
	}

	line := text[:crlfIdx]
	n = crlfIdx + 2

	if len(line) == 0 {
		return n, true, nil
	}

	colonIdx := strings.Index(line, ":")
	if colonIdx == -1 {
		return 0, false, fmt.Errorf("error: invalid header")
	}
	fieldName := line[:colonIdx]
	fieldValue := line[colonIdx+1:]
	if len(fieldName) == 0 || fieldName[len(fieldName)-1] == ' ' || fieldName[len(fieldName)-1] == '\t' {
		return 0, false, fmt.Errorf("error: invalid field name")
	}

	name := strings.ToLower(strings.TrimSpace(fieldName))
	for _, c := range name {
		if !('a' <= c && c <= 'z' || '0' <= c && c <= '9' || strings.Contains("!#$%&'*+-.^_`|~", string(c))) {
			return 0, false, fmt.Errorf("error: invalid character in field name: %v", string(c))
		}
	}
	value := strings.TrimSpace(fieldValue)

	if _, ok := h[strings.ToLower(name)]; ok {
		h[name] = h[name] + ", " + value
	} else {
		h[name] = value
	}

	return n, false, nil
}
