package main

import (
	"bytes"
	"log"

	"github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
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

	log.Fatal(coap.ListenAndServe("udp", ":1373", r))
}
