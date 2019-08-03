package httpbin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// WriteError :
func WriteError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"message": %q}`, error)
}

// Status : /status/{status}
func Status(w http.ResponseWriter, r *http.Request) {
	nodes := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	code, err := strconv.Atoi(nodes[len(nodes)-1])

	if err != nil {
		WriteError(w, err.Error(), 400)
		return
	}
	if code < 200 || code > 999 { // on status code, 1xx are also ok, but...
		WriteError(w, fmt.Sprintf("invalid WriteHeader code %v", code), 400)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	message := http.StatusText(code)
	fmt.Fprintf(w, `{"status": %d, "message": %q}`, code, message)
}
