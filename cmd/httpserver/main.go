package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/evanwiseman/httpfromtcp/internal/headers"
	"github.com/evanwiseman/httpfromtcp/internal/request"
	"github.com/evanwiseman/httpfromtcp/internal/response"
	"github.com/evanwiseman/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w *response.Writer, r *request.Request) {
	target := strings.TrimSpace(r.RequestLine.RequestTarget)
	if r.RequestLine.RequestTarget == "/myproblem" {
		handler400(w, r)
		return
	}
	if r.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, r)
		return
	}
	if strings.HasPrefix(target, "/httpbin") {
		handlerHttpbin(w, r)
		return
	}
	handler200(w, r)

}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOk)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>`)
	w.WriteHeaders(response.GetDefaultHeaders(len(body)))
	w.WriteBody(body)
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`)
	w.WriteHeaders(response.GetDefaultHeaders(len(body)))
	w.WriteBody(body)
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>`)
	w.WriteHeaders(response.GetDefaultHeaders(len(body)))
	w.WriteBody(body)
}

func handlerHttpbin(w *response.Writer, r *request.Request) {
	// Get the target
	endpoint := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin")
	url := "https://httpbin.org" + endpoint

	// Get the response from the url
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch from %s: %v", url, err)
		handler400(w, r)
		return
	}
	defer res.Body.Close()

	// Remove for chunked encoding
	res.Header.Del("content-length")

	header := headers.NewHeaders()
	header.Set("Transfer-Encoding", "Chunked")
	header.Set("Content-Type", res.Header.Get("Content-Type"))
	header.Set("Trailer", "X-Content-SHA256, X-Content-Length")
	w.WriteStatusLine(response.StatusOk)
	w.WriteHeaders(header)
	log.Print(header)
	length := 0
	buf := make([]byte, 1024)

	var body []byte
	for {
		n, err := res.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			handler400(w, r)
			return
		}
		length += n
		body = append(body, buf[:n]...)
		if _, err := w.WriteChunkedBody(buf[:n]); err != nil {
			handler400(w, r)
			return
		}
	}
	w.WriteChunkedBodyDone()

	hash := sha256.Sum256(body)

	trailer := headers.NewHeaders()
	trailer.Set("X-Content-SHA256", hex.EncodeToString(hash[:]))
	trailer.Set("X-Content-Length", fmt.Sprintf("%d", length))
	w.WriteTrailers(trailer)
}
