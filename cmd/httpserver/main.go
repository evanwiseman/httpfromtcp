package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		if err := w.Write(
			response.StatusBadRequest,
			[]byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)); err != nil {
			log.Printf("failed to write: %v", err)
		}
		return
	case "/myproblem":
		if err := w.Write(
			response.StatusInternalServerError,
			[]byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)); err != nil {
			log.Printf("failed to write: %v", err)
		}
		return
	case "/":
		if err := w.Write(
			response.StatusOk,
			[]byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)); err != nil {
			log.Printf("failed to write: %v", err)
		}
	}
}
