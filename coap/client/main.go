// Package main implements a simple CoAP (Constrained Application Protocol) client
// that connects to a CoAP server over UDP and performs GET requests.
//
// The client dials the server at localhost:1373 and sends a GET request to
// a specified path (default: "/b", or provided via command-line argument).
// It then reads and prints the response payload to the console.
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/plgd-dev/go-coap/v3/udp"
)

func main() {
	co, err := udp.Dial("localhost:1373")
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	defer func() {
		_ = co.Close()
	}()

	path := "/b"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Printf("Error sending request: %s\n", err)

		return
	}

	b, err := resp.ReadBody()
	if err != nil {
		log.Printf("Error reading the response body: %s\n", err)

		return
	}

	log.Printf("Response payload: %s", string(b))
}
