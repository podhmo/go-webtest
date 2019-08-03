package client

import (
	"io"
	"net/http"
	"net/url"
	"strings"

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
	return c.Internal.Request("GET", path, nil)
}

// Head :
func (c *Adapter) Head(path string) (response.Response, error, func()) {
	return c.Internal.Request("HEAD", path, nil)
}

// Post :
func (c *Adapter) Post(
	path, contentType string,
	body io.Reader,
) (response.Response, error, func()) {
	return c.Internal.Request("POST", path, body, func(req *http.Request) {
		req.Header.Set("Content-Type", contentType)
	})
}

// PostForm :
func (c *Adapter) PostForm(
	path string,
	data url.Values,
) (response.Response, error, func()) {
	return c.Post(
		path,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
}

// PostJSON :
func (c *Adapter) PostJSON(
	path string,
	body io.Reader,
) (response.Response, error, func()) {
	return c.Post(path, "application/json", body)
}

// Internal :
type Internal interface {
	Do(req *http.Request) (response.Response, error, func())
	Request(method string, path string, body io.Reader, options ...func(*http.Request)) (response.Response, error, func())
}
