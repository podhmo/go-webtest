package tripperware

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestTripperware(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, req.URL.Query().Encode())
	}

	cases := []struct {
		msg   string
		want  string
		path  string
		stack List
	}{
		{
			msg:  "zero decorator",
			want: "x=1",
			path: "/?x=1",
		},
		{
			msg:  "one decorator",
			want: "x=1&y=2",
			path: "/?x=1",
			stack: Stack(
				func(next http.RoundTripper) http.RoundTripper {
					return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						q := req.URL.Query()
						q.Add("y", "2")
						req.URL.RawQuery = q.Encode()
						return next.RoundTrip(req)
					})
				},
			),
		},
		{
			msg:  "two decorator",
			want: "x=1&y=2&z=3",
			path: "/?x=1",
			stack: Stack(
				func(next http.RoundTripper) http.RoundTripper {
					return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						q := req.URL.Query()
						q.Add("y", "2")
						req.URL.RawQuery = q.Encode()
						return next.RoundTrip(req)
					})
				},
				func(next http.RoundTripper) http.RoundTripper {
					return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
						q := req.URL.Query()
						q.Add("z", "3")
						req.URL.RawQuery = q.Encode()
						return next.RoundTrip(req)
					})
				},
			),
		},
	}

	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(handler))

			client := c.stack.DecorateClient(&http.Client{}, true)
			res, err := client.Get(fmt.Sprintf("%s%s", ts.URL, c.path))
			if err != nil {
				t.Fatal(err)
			}

			var b strings.Builder
			if _, err := io.Copy(&b, res.Body); err != nil {
				t.Fatal(err)
			}
			got := b.String()
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("want %s, but %s", c.want, got)
			}
		})
	}
}
