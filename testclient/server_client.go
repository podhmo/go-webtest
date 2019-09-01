package testclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

type roundTripperWrapper struct {
	http.RoundTripper
	Wrap func(inner http.RoundTripper) http.RoundTripper
}

// TODO: logging interface

// ServerClient :
type ServerClient struct {
	Client    *http.Client
	Transport http.RoundTripper

	Server   *httptest.Server
	BasePath string // need?
}

// Do :
func (c *ServerClient) Do(req *http.Request) (Response, error, func()) {
	client := GetInternalClientWith(c.Client, c.Transport)

	var adapter *ResponseAdapter
	var raw *http.Response
	var once sync.Once

	raw, err := client.Do(req)
	if err != nil {
		return nil, err, nil
	}

	adapter = NewResponseAdapter(
		func() *http.Response {
			once.Do(func() {
				adapter.AddTeardown(raw.Body.Close)
			})
			return raw
		},
	)
	return adapter, err, adapter.Close
}

// NewRequest :
func (c *ServerClient) NewRequest(
	method string,
	path string,
	body io.Reader,
) (*http.Request, error) {
	url := internal.URLJoin(c.Server.URL, internal.URLJoin(c.BasePath, path))
	return http.NewRequest(method, url, body)
}
