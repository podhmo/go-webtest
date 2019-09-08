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
			return nil, &statusError{want: code, response: res, text: b.String()}
		})
	}
}

type statusError struct {
	want int

	response *http.Response
	text     string
}

func (err *statusError) Error() string {
	return fmt.Sprintf(
		"status code, got is \"%[1]d %[2]s\", but want is \"%[3]d %[4]s\"\n response is %[5]s",
		err.response.StatusCode,
		http.StatusText(err.response.StatusCode),
		err.want,
		http.StatusText(err.want),
		err.text,
	)
}
