package testclient

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestRoundTripperDecorator(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, req.URL.Query().Encode())
	}

	cases := []struct {
		msg       string
		want      string
		path      string
		decorator RoundTripperDecorator
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
			decorator: FuncRoundTripper{
				Fn: func(inner http.RoundTripper, req *http.Request) (*http.Response, error) {
					q := req.URL.Query()
					q.Add("y", "2")
					req.URL.RawQuery = q.Encode()
					return inner.RoundTrip(req)
				},
			},
		},
		{
			msg:  "two decorator",
			want: "x=1&y=2&z=3",
			path: "/?x=1",
			decorator: FuncRoundTripper{
				Fn: func(inner http.RoundTripper, req *http.Request) (*http.Response, error) {
					q := req.URL.Query()
					q.Add("y", "2")
					req.URL.RawQuery = q.Encode()
					return inner.RoundTrip(req)
				},
			}.Decorate(FuncRoundTripper{
				Fn: func(inner http.RoundTripper, req *http.Request) (*http.Response, error) {
					q := req.URL.Query()
					q.Add("z", "3")
					req.URL.RawQuery = q.Encode()
					return inner.RoundTrip(req)
				},
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(handler))

			client := &http.Client{}
			if c.decorator != nil {
				client.Transport = c.decorator.Decorate(client.Transport)
			}

			res, err := client.Get(fmt.Sprintf("%s%s", ts.URL, c.path))
			if err != nil {
				t.Fatal(err)
			}

			var b strings.Builder
			io.Copy(&b, res.Body)
			got := b.String()
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("want %s, but %s", c.want, got)
			}
		})
	}
}
