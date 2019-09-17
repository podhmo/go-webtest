package tripperware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
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

			b, err := httputil.DumpResponse(res, true)
			if err != nil {
				return nil, err
			}
			return nil, &statusError{want: code, response: res, text: string(b)}
		})
	}
}

type statusError struct {
	want int

	response *http.Response
	text     string
}

func (err *statusError) Error() string {
	var b strings.Builder
	w := &b

	fmt.Fprintln(w, "\x1b[32mResponse: ------------------------------\x1b[0m")
	fmt.Fprint(w, err.text)
	if !strings.HasSuffix(strings.TrimRight(err.text, "  \t"), "\n") {
		fmt.Fprintln(w, "")
	}
	fmt.Fprintln(w, "\x1b[32m----------------------------------------\x1b[0m")

	return fmt.Sprintf(
		"status code, got is \"%[1]d %[2]s\", but want is \"%[3]d %[4]s\"\n%[5]s",
		err.response.StatusCode,
		http.StatusText(err.response.StatusCode),
		err.want,
		http.StatusText(err.want),
		w.String(),
	)
}
