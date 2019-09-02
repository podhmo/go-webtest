package testclient

import (
	"net/http"
	"sync"

	"github.com/podhmo/go-webtest/testclient/internal"
	"github.com/podhmo/go-webtest/tripperware"
)

// RealClient :
type RealClient struct {
	URL    string
	Client *http.Client
}

// Do :
func (c *RealClient) Do(
	req *http.Request,
	config *Config,
) (Response, error) {
	stack := tripperware.Stack(append(defaultTripperwares, config.Tripperwares...)...)
	cloned := true
	client := stack.DecorateClient(c.Client, cloned)

	var adapter *ResponseAdapter
	var raw *http.Response
	var once sync.Once

	raw, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	adapter = NewResponseAdapter(
		func() *http.Response {
			once.Do(func() {
				adapter.AddTeardown(raw.Body.Close)
			})
			return raw
		},
	)
	return adapter, err
}

// NewRequest :
func (c *RealClient) NewRequest(
	method string,
	path string,
	config *Config,
) (*http.Request, error) {
	url := internal.URLJoin(c.URL, internal.URLJoin(config.BasePath, path))
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
