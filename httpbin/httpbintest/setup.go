package httpbintest

import (
	"net/http"
	"net/http/httptest"

	"github.com/podhmo/go-webtest/httpbin"
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
