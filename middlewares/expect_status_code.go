package middlewares

import (
	"fmt"
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
)

// ExpectCode :
func ExpectCode(code int) func(*webtest.Config) {
	return func(c *webtest.Config) {
		c.Middlewares = append(c.Middlewares, NewMiddleware(func(
			t testing.TB,
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
		err.response.LazyBodyString(),
	)
}
