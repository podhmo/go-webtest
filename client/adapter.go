package client

import (
	"io"
	"net/http"

	"github.com/podhmo/go-webtest/client/response"
)

// Adapter :
type Adapter struct {
	Internal Internal
}

// Do :
func (c *Adapter) Do(req *http.Request) (response.Response, error, func()) {
	return c.Internal.Do(req)
}

// Get :
func (c *Adapter) Get(path string) (response.Response, error, func()) {
	var body io.Reader // xxx (TODO: functional options)
	return c.Internal.Request("GET", path, body)
}

// Internal :
type Internal interface {
	Do(req *http.Request) (response.Response, error, func())
	Request(method string, path string, body io.Reader) (response.Response, error, func())
}
