package main

import (
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	addr = "127.0.0.1:9994"
)

func main() {
	fmt.Printf("Destination server listening for connections on %s...\n", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to create listener with %v\n", err)
		os.Exit(1)
	}

	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		}),
	}
	log.Fatal(server.Serve(listener))
}
