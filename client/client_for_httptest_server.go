package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

// HTTPTestServerClient :
type HTTPTestServerClient struct {
	client   *http.Client
	Server   *httptest.Server
	BasePath string // need?
}

// Do :
func (c *HTTPTestServerClient) Do(req *http.Request) (Response, error, func()) {
	client := c.client
	if c.client == nil {
		client = http.DefaultClient
	}

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
func (c *HTTPTestServerClient) NewRequest(
	method string,
	path string,
	body io.Reader,
) (*http.Request, error) {
	url := internal.URLJoin(c.Server.URL, internal.URLJoin(c.BasePath, path))
	return http.NewRequest(method, url, body)
}
