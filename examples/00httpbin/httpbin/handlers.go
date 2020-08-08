package httpbin

import (
	"encoding/json"
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

// Get : /get
func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	d := map[string]interface{}{
		"args":    r.URL.Query(),
		"headers": r.Header,
		"url":     r.URL.String(),
	}

	if err := encoder.Encode(d); err != nil {
		panic(err)
	}
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

// BasicAuth : /basic-auth/{user}/{passwd} Prompts the user for authorization using HTTP Basic Auth.
func BasicAuth(w http.ResponseWriter, req *http.Request) {
	nodes := strings.Split(strings.TrimSuffix(req.URL.Path, "/"), "/")
	basicAuthUser := nodes[len(nodes)-2]
	basicAuthPassword := nodes[len(nodes)-1]

	user, pass, ok := req.BasicAuth()
	if !ok || user != basicAuthUser || pass != basicAuthPassword {
		w.Header().Add("WWW-Authenticate", `Basic realm="Fake Realm"`)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"authenticated": false}`)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, `{"authenticated": true, "user": %q}`, user)
}
