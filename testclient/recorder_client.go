package testclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

// RecorderClient :
type RecorderClient struct {
	Handler  http.Handler
	BasePath string

	Transport http.RoundTripper
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
) (Response, error, func()) {
	var adapter *ResponseAdapter
	var once sync.Once

	transport := getWrappedTransport(c, c.Transport)
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
	body io.Reader,
) (*http.Request, error) {
	url := internal.URLJoin(c.BasePath, path)
	return httptest.NewRequest(method, url, body), nil
}
