package hook

import (
	"net/http"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/testclient"
)

// Response :
type Response = testclient.Response

// Hook :
type Hook = webtest.Hook

// NewHook :
func NewHook(wrap func(res Response, req *http.Request) error) Hook {
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
