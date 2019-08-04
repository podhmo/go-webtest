package middlewares

import (
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/client"
)

// Response :
type Response = client.Response

// Middleware :
type Middleware = webtest.Middleware

// NewMiddleware :
func NewMiddleware(wrap func(t testing.TB, res Response, req *http.Request) error) Middleware {
	return func(
		t testing.TB,
		req *http.Request,
		inner func(*http.Request) (Response, error, func()),
	) (Response, error, func()) {
		res, err, teardown := inner(req)
		if err != nil {
			return res, err, teardown
		}
		if err := wrap(t, res, req); err != nil {
			return res, err, teardown
		}
		return res, err, teardown
	}
}
