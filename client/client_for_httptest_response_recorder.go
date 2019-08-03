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

// Request :
func (c *HTTPTestResponseRecorderClient) Request(
	method string,
	path string,
	body io.Reader,
	options ...func(*http.Request),
) (Response, error, func()) {
	url := internal.URLJoin(c.BasePath, path)
	req := httptest.NewRequest(method, url, body)
	for _, opt := range options {
		opt(req)
	}
	return c.Do(req)
}
