package middlewares

import (
	"net/http"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/client"
)

// Response :
type Response = client.Response

// Middleware :
type Middleware = webtest.Middleware

// NewMiddleware :
func NewMiddleware(wrap func(res Response, req *http.Request) error) Middleware {
	return func(
		req *http.Request,
		inner func(*http.Request) (Response, error, func()),
	) (Response, error, func()) {
		res, err, teardown := inner(req)
		if err != nil {
			return res, err, teardown
		}
		if err := wrap(res, req); err != nil {
			return res, err, teardown
		}
		return res, err, teardown
	}
}
