package testclient

import (
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

// RecorderClient :
type RecorderClient struct {
	Handler http.Handler
}

// RoundTrip :
func (c *RecorderClient) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO: accessing headder information
	w := httptest.NewRecorder()
	c.Handler.ServeHTTP(w, req)
	return w.Result(), nil
}

// Do :
func (c *RecorderClient) Do(
	req *http.Request,
	config *Config,
) (Response, error, func()) {
	var adapter *ResponseAdapter
	var once sync.Once

	transport := getDecoratepedTransport(c, config.Decorator)
	res, err := transport.RoundTrip(req)

	adapter = NewResponseAdapter(
		func() *http.Response {
			once.Do(func() {
				adapter.AddTeardown(res.Body.Close)
			})
			return res
		},
	)
	return adapter, err, adapter.Close
}

// NewRequest :
func (c *RecorderClient) NewRequest(
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
