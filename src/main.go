package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8080"
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	http.Handle("/strings/hello", HelloWorldHandler{})
	http.Handle("/strings/async-hello", AsyncHelloWorldHandler{})
	http.Handle("/static/", StaticFileServer{})

	log.Printf("Serving on port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
