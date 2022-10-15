package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/andybalholm/brotli"
)

const maxNameLength = 100

type HelloWorldHandler struct{}

func (h HelloWorldHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()

	name := qParams.Get("name")
	if name == "" {
		name = "world"
	} else if len(name) > maxNameLength {
		http.Error(w, fmt.Sprintf("Name must be <= %v characters", maxNameLength), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain")

	compressor := brotli.HTTPCompressor(w, req)
	compressor.Write([]byte("Hello, " + name + "!"))

	err := compressor.Close()
	if err != nil {
		log.Printf("Failed to close compressor: %v", err.Error())
		http.Error(w, "Compression failure", http.StatusInternalServerError)
		return
	}
}
