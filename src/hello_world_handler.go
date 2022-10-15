package main

import "net/http"

type HelloWorldHandler struct{}

func (h HelloWorldHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()

	name := qParams.Get("name")
	if name == "" {
		name = "world"
	}

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("Hello, " + name + "!"))
}
