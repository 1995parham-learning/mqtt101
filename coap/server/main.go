package main

import (
	"bytes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/plgd-dev/go-coap/v3/udp"
	coapNet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"github.com/plgd-dev/go-coap/v3/options"
)

func handleA(w mux.ResponseWriter, _ *mux.Message) {
	if err := w.SetResponse(codes.GET, message.TextPlain, bytes.NewReader([]byte("hello world"))); err != nil {
		log.Printf("cannot set response: %v", err)
	}
}

func main() {
	r := mux.NewRouter()

	if err := r.Handle("/a", mux.HandlerFunc(handleA)); err != nil {
		log.Fatal(err)
	}

	udpConn, err := coapNet.NewListenUDP("udp", ":1373")
	if err != nil {
		log.Fatal(err)
	}
	defer udpConn.Close()

	server := udp.NewServer(options.WithMux(r))

	go func() {
		if err := server.Serve(udpConn); err != nil {
			log.Printf("CoAP server error: %v", err)
		}
	}()

	log.Println("CoAP server listening on :1373")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("signal caught - shutting down")

	server.Stop()

	log.Println("shutdown complete")
}
