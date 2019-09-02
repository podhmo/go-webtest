package hook

import (
	"fmt"
	"net/http"
	"testing"
)

// ExpectCode :
func ExpectCode(t testing.TB, code int) Hook {
	return NewHook(func(res Response, req *http.Request) error {
		if res.Code() != code {
			return &statusError{code: code, response: res}
		}
		return nil
	})
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
