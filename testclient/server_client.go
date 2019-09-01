package testclient

import (
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

// ServerClient :
type ServerClient struct {
	Server *httptest.Server
	Client *http.Client
}

// Do :
func (c *ServerClient) Do(
	req *http.Request,
	config *Config,
) (Response, error, func()) {
	client := getInternalClient(c.Client, config.Decorator)

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
	config *Config,
) (*http.Request, error) {
	url := internal.URLJoin(c.Server.URL, internal.URLJoin(config.BasePath, path))
	req, err := http.NewRequest(method, url, config.Body)
	if err != nil {
		return req, err
	}

	if config.Query != nil {
		q := config.Query
		for k, vs := range req.URL.Query() {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}
