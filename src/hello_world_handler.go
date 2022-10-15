package main

import (
	"log"
	"net/http"

	"github.com/andybalholm/brotli"
)

type HelloWorldHandler struct{}

func (h HelloWorldHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()

	name := qParams.Get("name")
	if name == "" {
		name = "world"
	}

	w.Header().Add("Content-Type", "text/plain")

	compressor := brotli.HTTPCompressor(w, req)
	compressor.Write([]byte("Hello, " + name + "!"))

	err := compressor.Close()
	if err != nil {
		log.Printf("Failed to close compressor: %v", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
