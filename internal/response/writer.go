package response

import (
	"fmt"
	"io"
	"net"

	"github.com/evanwiseman/httpfromtcp/internal/headers"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{
		writer: conn,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if err := WriteStatusLine(w.writer, statusCode); err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if err := WriteHeaders(w.writer, headers); err != nil {
		return err
	}

	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	numBytes := len(p)
	return w.WriteBody([]byte(fmt.Sprintf("%X\r\n%s\r\n", numBytes, p)))
}

func (w *Writer) WriteChunkedBodyDone() error {
	_, err := w.writer.Write([]byte("0\r\n"))
	return err
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	return w.WriteHeaders(h)
}
