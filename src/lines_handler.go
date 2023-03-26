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

	nStr := qParams.Get("n")
	n := 0
	if nStr != "" {
		n64, err := strconv.ParseInt(qParams.Get("n"), 10, strconv.IntSize)
		if err != nil {
			http.Error(w, "Could not parse param n as an int", http.StatusBadRequest)
			return
		}
		n = int(n64)
	}

	tags := strings.Builder{}
	tags.WriteString("<ol>\n")
	for i := 1; i <= n; i++ {
		tags.WriteString(fmt.Sprintf("  <li>Item number: %v</li>\n", i))
	}
	tags.WriteString("</ol>")
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	MaybeCompress(w, req, []byte(tags.String()))
}
