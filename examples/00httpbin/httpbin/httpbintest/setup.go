package httpbintest

import (
	"net/http"
	"net/http/httptest"

	"m/httpbin"
)

// NewTestAPIServer :
func NewTestAPIServer() (*httptest.Server, func()) {
	ts := httptest.NewServer(httpbin.Handler())
	return ts, ts.Close
}

// NewTestHandler :
func NewTestHandler() http.HandlerFunc {
	return httpbin.Handler().ServeHTTP
}
