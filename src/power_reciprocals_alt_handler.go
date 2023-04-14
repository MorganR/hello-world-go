package main

import (
	"net/http"
	"strconv"
)

// PowerReciprocalsAltHandler computes a convergent series with "n" terms.
type PowerReciprocalsAltHandler struct{}

func (h PowerReciprocalsAltHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()

	nStr := qParams.Get("n")
	n := int64(0)
	if nStr != "" {
		var err error
		n, err = strconv.ParseInt(qParams.Get("n"), 10, strconv.IntSize)
		if err != nil {
			http.Error(w, "Could not parse param n as an int", http.StatusBadRequest)
			return
		}
	}

	result := 0.0
	power := 0.5
	for i := int64(0); i < n; i++ {
		power *= 2
		result += 1 / power

		i++

		if i < n {
			power *= 2
			result -= 1 / power
		}
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(strconv.FormatFloat(result, 'f', -1, 64)))
}
