package tripperware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// ExpectCode :
func ExpectCode(t testing.TB, code int) Ware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			res, err := next.RoundTrip(req)
			if err != nil {
				return nil, err
			}
			if res.StatusCode == code {
				return res, nil
			}

			defer res.Body.Close()
			var b bytes.Buffer
			if _, err := io.Copy(&b, res.Body); err != nil {
				return nil, err
			}
			return nil, &statusError{expected: code, response: res, text: b.String()}
		})
	}
}

type statusError struct {
	expected int

	response *http.Response
	text     string
}

func (err *statusError) Error() string {
	return fmt.Sprintf(
		"status code, expected %d, but actual %d\n response: %s",
		err.expected,
		err.response.StatusCode,
		err.text,
	)
}
