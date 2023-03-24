package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-ocf/go-coap"
)

func main() {
	co, err := coap.Dial("udp", "localhost:1373")

	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	path := "/b"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := co.GetWithContext(ctx, path)

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	log.Printf("Response payload: %v", resp.Payload())
}
