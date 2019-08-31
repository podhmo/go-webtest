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
}

// Do :
func (c *RecorderClient) Do(
	req *http.Request,
) (Response, error, func()) {
	var adapter *ResponseAdapter
	var res *http.Response
	var once sync.Once

	w := httptest.NewRecorder()
	c.Handler.ServeHTTP(w, req)

	adapter = NewResponseAdapter(
		func() *http.Response {
			once.Do(func() {
				res = w.Result()
				adapter.AddTeardown(res.Body.Close)
			})
			return res
		},
	)
	return adapter, nil, adapter.Close
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
