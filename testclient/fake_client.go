package testclient

import (
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/testclient/internal"
	"github.com/podhmo/go-webtest/tripperware"
)

// FakeClient :
type FakeClient struct {
	Handler http.Handler
}

// RoundTrip :
func (c *FakeClient) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO: accessing headder information
	w := httptest.NewRecorder()
	c.Handler.ServeHTTP(w, req)
	return w.Result(), nil
}

// Do :
func (c *FakeClient) Do(
	req *http.Request,
	config *Config,
) (Response, error) {
	var adapter *ResponseAdapter
	var once sync.Once

	stack := tripperware.Stack(append(defaultTripperwares, config.Tripperwares...)...)
	res, err := stack.DecorateRoundTripper(c).RoundTrip(req)
	if err != nil {
		return nil, err
	}
	res.Request = req

	adapter = NewResponseAdapter(
		func() *http.Response {
			once.Do(func() {
				adapter.AddTeardown(res.Body.Close)
			})
			return res
		},
	)
	return adapter, err
}

// NewRequest :
func (c *FakeClient) NewRequest(
	method string,
	path string,
	config *Config,
) (*http.Request, error) {
	url := internal.URLJoin(config.BasePath, path)
	req := httptest.NewRequest(method, url, config.Body)

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
