package response

import (
	"fmt"
	"net"

	"github.com/evanwiseman/httpfromtcp/internal/headers"
)

type Writer struct {
	conn net.Conn
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		conn: conn,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if err := WriteStatusLine(w.conn, statusCode); err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if err := WriteHeaders(w.conn, headers); err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) error {
	_, err := w.conn.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) Write(statusCode StatusCode, body []byte) error {
	if err := w.WriteStatusLine(statusCode); err != nil {
		return fmt.Errorf("error: failed to write status line: %w", err)
	}

	headers := GetDefaultHeaders(len(body))
	if err := w.WriteHeaders(headers); err != nil {
		return fmt.Errorf("error: failed to write headers: %w", err)
	}

	if err := w.WriteBody(body); err != nil {
		return fmt.Errorf("error: failed to write body: %w", err)
	}
	return nil
}
