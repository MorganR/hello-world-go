package main

import (
	"fmt"
	"net/http"
)

const maxNameLength = 500

// HelloWorldHandler provides a greeting, using the optional "name" query parameter.
type HelloWorldHandler struct{}

func (h HelloWorldHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()

	var greeting []byte
	name := qParams.Get("name")
	if name == "" {
		greeting = []byte("Hello, world!")
	} else if len(name) > maxNameLength {
		http.Error(w, fmt.Sprintf("Name must be <= %v characters", maxNameLength), http.StatusBadRequest)
		return
	} else {
		greeting = []byte("Hello, " + name + "!")
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	MaybeCompress(w, req, greeting)
}
