package main

import (
	"net/http"
	"time"
)

// AsyncHelloWorldHandler provides a greeting after a delay.
type AsyncHelloWorldHandler struct{}

var response = []byte("Hello, world!")

func (h AsyncHelloWorldHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Millisecond * 15)
	w.Write(response)
}
