package hook

import (
	"fmt"
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
)

// ExpectCode :
func ExpectCode(t testing.TB, code int) webtest.Option {
	return func(c *webtest.Config) {
		c.Hooks = append(c.Hooks, NewHook(func(
			res Response,
			req *http.Request,
		) error {
			if res.Code() != code {
				return &statusError{code: code, response: res}
			}
			return nil
		}))
	}
}

type statusError struct {
	code     int
	response Response
}

func (err *statusError) Error() string {
	return fmt.Sprintf(
		"status code, expected %d, but actual %d\n response: %s",
		err.code,
		err.response.Code(),
		err.response.LazyText(),
	)
}
