package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.sock")
		return
	}
	sockPath := os.Args[1]

	fmt.Println("Proxy Server listening on Unix Domain Socket")

	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		fmt.Printf("Failed to create listener %s with %v\n", sockPath, err)
		os.Exit(1)
	}

	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handleHTTP(w, r)
		}),
	}

	server.Serve(listener)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(string(requestDump))
	}

	u, err := url.Parse("http://localhost:9994/")
	if err != nil {
		log.Fatal(err)
	}
	req.URL = u

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
