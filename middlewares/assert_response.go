package middlewares

import (
	"fmt"
	"net/http"
)

// AssertResponse
func AssertResponse(code int) Middleware {
	return NewMiddleware(func(res Response, req *http.Request) error {
		if res.StatusCode() != code {
			return &assertError{code: code, response: res}
		}
		return nil
	})
}

type assertError struct {
	code     int
	response Response
}

func (err *assertError) Error() string {
	return fmt.Sprintf(
		"status code, expected %d, but actual %d\n response: %s",
		err.code,
		err.response.StatusCode(),
		err.response.LazyBodyString(),
	)
}
