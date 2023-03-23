package main

import (
	"net/http"
	"strconv"
)

// PowerReciprocalsAltHandler computes a convergent series.
type PowerReciprocalsAltHandler struct{}

func (h PowerReciprocalsAltHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	qParams := req.URL.Query()

	n, err := strconv.ParseInt(qParams.Get("n"), 10, 64)
	if err != nil {
		http.Error(w, "Could not parse param n as an int", http.StatusBadRequest)
		return
	}

	result := float64(0.0)
	power := float64(0.5)
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
