package middlewares

import (
	"fmt"
	"net/http"

	webtest "github.com/podhmo/go-webtest"
)

// ExpectStatusCode :
func ExpectStatusCode(code int) func(*webtest.Config) {
	return func(c *webtest.Config) {
		c.Middlewares = append(c.Middlewares, NewMiddleware(func(res Response, req *http.Request) error {
			if res.StatusCode() != code {
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
		err.response.StatusCode(),
		err.response.LazyBodyString(),
	)
}
