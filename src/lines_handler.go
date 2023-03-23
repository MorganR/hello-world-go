package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// LinesHandler prints a formatted HTML content with "n" list items.
type LinesHandler struct{}

func (h LinesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()

	n64, err := strconv.ParseInt(qParams.Get("n"), 10, strconv.IntSize)
	if err != nil {
		http.Error(w, "Could not parse param n as an int", http.StatusBadRequest)
		return
	}
	n := int(n64)

	tags := make([]string, 0, n+2)
	tags = append(tags, "<ol>")
	for i := 1; i <= n; i++ {
		tags = append(tags, fmt.Sprintf("  <li>Item number: %v</li>", i))
	}
	tags = append(tags, "</ol>")
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	MaybeCompress(w, req, []byte(strings.Join(tags, "\n")))
}
