package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

// HTTPTestResponseRecorderClient :
type HTTPTestResponseRecorderClient struct {
	HandlerFunc http.HandlerFunc
	BasePath    string
}

// Do :
func (c *HTTPTestResponseRecorderClient) Do(
	req *http.Request,
) (Response, error, func()) {
	var adapter *ResponseAdapter
	var raw *http.Response
	var once sync.Once

	w := httptest.NewRecorder()
	c.HandlerFunc(w, req)

	adapter = NewResponseAdapter(
		func() *http.Response {
			once.Do(func() {
				raw = w.Result()
				adapter.AddTeardown(raw.Body.Close)
			})
			return raw
		},
	)
	return adapter, nil, adapter.Close
}

// NewRequest :
func (c *HTTPTestResponseRecorderClient) NewRequest(
	method string,
	path string,
	body io.Reader,
) (*http.Request, error) {
	url := internal.URLJoin(c.BasePath, path)
	return httptest.NewRequest(method, url, body), nil
}
