package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/internal"
	"github.com/podhmo/go-webtest/client/response"
)

// HTTPTestServerClient :
type HTTPTestServerClient struct {
	client   *http.Client
	Server   *httptest.Server
	BasePath string // need?
}

// Do :
func (c *HTTPTestServerClient) Do(req *http.Request) (response.Response, error, func()) {
	client := c.client
	if c.client == nil {
		client = http.DefaultClient
	}

	var adapter *response.Adapter
	var raw *http.Response
	var once sync.Once

	raw, err := client.Do(req)
	if err != nil {
		return nil, err, nil
	}

	adapter = response.NewAdapter(
		func() *http.Response {
			once.Do(func() {
				adapter.AddTeardown(raw.Body.Close)
			})
			return raw
		},
	)
	return adapter, err, adapter.Close
}

// Get :
func (c *HTTPTestServerClient) Get(path string) (response.Response, error, func()) {
	url := internal.URLJoin(c.Server.URL, internal.URLJoin(c.BasePath, path))
	var body io.Reader // xxx (TODO: functional options)
	req, err := http.NewRequest("GET", url, body)
	if err != nil {
		return nil, err, nil
	}
	return c.Do(req)
}

// TODO: more methods
