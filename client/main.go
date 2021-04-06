package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	protocol = "unix"
)

func main() {
	post := flag.String("d", "", "data to POST")
	help := flag.Bool("h", false, "usage help")
	sockPath := flag.String("p", "/tmp/echo.sock", "path to unix socket")
	uri := flag.String("u", "/", "HTTP URI")
	flag.Parse()

	if *help || len(flag.Args()) > 3 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "[-d data] [-p path] [-u uri]")
		flag.PrintDefaults()
		os.Exit(0)
	}

	var (
		proxyFinished, clientFinished chan bool
	)

	go runProxy(proxyFinished, *sockPath)
	go connectProxy(clientFinished, *post, *sockPath, *uri)

	<-proxyFinished
	<-clientFinished
}

func runProxy(finished chan bool, sockPath string) {
	binary, err := exec.LookPath("./bin/proxy")
	if err != nil {
		fmt.Printf("Could not find %s: %v\n", binary, err)
		os.Exit(1)
	}

	cmd := exec.Command(binary, sockPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Execution of %s failed with %v\n", binary, err)
		os.Exit(1)
	}

	finished <- true
}

func connectProxy(finished chan bool, post, sockPath, uri string) {
	var (
		response *http.Response
		err      error
	)

	client := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sockPath)
			},
		},
	}

	delay := 100
	num := 3

	for i := 0; i <= num; i++ {
		backoff := backoff(delay, i)
		fmt.Printf("delay: %v\n", backoff)
		time.Sleep(backoff)
		if len(post) == 0 {
			response, err = client.Get("http://unix" + uri)
		} else {
			response, err = client.Post("http://unix"+uri, "application/octet-stream", strings.NewReader(post))
		}

		if err != nil {
			fmt.Printf("Error connecting to the socket %s: %v\n", "http://unix"+uri, err)
			if i == num {
				os.Exit(1)
			}
		}
	}

	// copy to stdout
	io.Copy(os.Stdout, response.Body)

	finished <- true
}

func backoff(delay, n int) time.Duration {
	return time.Duration(float64(delay)*math.Pow(float64(2), float64(n))) * time.Millisecond
}
