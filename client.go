package webtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/podhmo/go-webtest/client"
	"github.com/podhmo/go-webtest/client/response"
)

// Internal :
type Internal interface {
	Do(req *http.Request) (response.Response, error, func())
	Request(method string, path string, body io.Reader, options ...func(*http.Request)) (response.Response, error, func())
}

// Client :
type Client struct {
	Internal Internal
}

// Do :
func (c *Client) Do(req *http.Request) (response.Response, error, func()) {
	return c.Internal.Do(req)
}

// Get :
func (c *Client) Get(path string) (response.Response, error, func()) {
	return c.Internal.Request("GET", path, nil)
}

// Head :
func (c *Client) Head(path string) (response.Response, error, func()) {
	return c.Internal.Request("HEAD", path, nil)
}

// Post :
func (c *Client) Post(
	path, contentType string,
	body io.Reader,
) (response.Response, error, func()) {
	return c.Internal.Request("POST", path, body, func(req *http.Request) {
		req.Header.Set("Content-Type", contentType)
	})
}

// PostForm :
func (c *Client) PostForm(
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
func (c *Client) PostJSON(
	path string,
	body io.Reader,
) (response.Response, error, func()) {
	return c.Post(path, "application/json", body)
}

// Response :
type Response = response.Response

// NewClientFromTestServer :
func NewClientFromTestServer(ts *httptest.Server, options ...func(*Config)) *Client {
	c := &Config{}
	for _, opt := range options {
		opt(c)
	}
	return &Client{
		Internal: &client.HTTPTestServerClient{
			Server:   ts,
			BasePath: c.BasePath,
		},
	}
}

// NewClientFromHandler :
func NewClientFromHandler(handlerFunc http.HandlerFunc, options ...func(*Config)) *Client {
	c := &Config{}
	for _, opt := range options {
		opt(c)
	}
	return &Client{
		Internal: &client.HTTPTestResponseRecorderClient{
			HandlerFunc: handlerFunc,
			BasePath:    c.BasePath,
		},
	}
}

// WithBasePath :
func WithBasePath(basePath string) func(*Config) {
	return func(c *Config) {
		c.BasePath = basePath
	}
}

// Config :
type Config struct {
	BasePath string
}
