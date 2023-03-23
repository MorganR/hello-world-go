package main

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

const minCompressionLength = 256
const gzipCompressionLevel = 6

func MaybeCompress(w http.ResponseWriter, req *http.Request, response []byte) {
	if len(response) < minCompressionLength ||
		!acceptsGzip(req) {
		w.Write(response)
		return
	}

	w.Header().Add("Content-Encoding", "gzip")
	compressor, err := gzip.NewWriterLevel(w, gzipCompressionLevel)
	if err != nil {
		log.Printf("Failed to create compressor: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	compressor.Write(response)

	err = compressor.Close()
	if err != nil {
		log.Printf("Failed to close compressor: %v", err.Error())
		http.Error(w, "Compression failure", http.StatusInternalServerError)
		return
	}
}

func acceptsGzip(req *http.Request) bool {
	acceptEncodings := strings.Split(req.Header.Get("Accept-Encoding"), ",")
	for _, enc := range acceptEncodings {
		if strings.TrimSpace(enc) == "gzip" {
			return true
		}
	}
	return false
}
