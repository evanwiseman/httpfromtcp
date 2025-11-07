package response

import (
	"fmt"
	"io"
)

type StatusCode int

const (
	StatusOk                  = 200
	StatusBadRequest          = 400
	StatusInternalServerError = 500
)

func GetStatusLine(statusCode StatusCode) []byte {
	var reason string
	switch statusCode {
	case StatusOk:
		reason = "OK"
	case StatusBadRequest:
		reason = "BAD REQUEST"
	case StatusInternalServerError:
		reason = "INTERNAL SERVER ERROR"
	}

	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reason))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

	_, err := w.Write(GetStatusLine(statusCode))
	if err != nil {
		return err
	}
	return nil
}
