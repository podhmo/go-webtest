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
	return Hook(func(
		req *http.Request,
		inner func(*http.Request) (Response, error),
	) (res Response, err error) {
		res, err = inner(req)
		if err != nil {
			return
		}
		err = wrap(res, req)
		return
	})
}
