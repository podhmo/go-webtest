package webtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/podhmo/go-webtest/client"
)

// Response :
type Response = client.Response

// Internal :
type Internal interface {
	Do(req *http.Request) (Response, error, func())
	Request(method string, path string, body io.Reader, options ...func(*http.Request)) (Response, error, func())
}

// Client :
type Client struct {
	Internal Internal
}

// Do :
func (c *Client) Do(req *http.Request) (Response, error, func()) {
	return c.Internal.Do(req)
}

// Get :
func (c *Client) Get(path string) (Response, error, func()) {
	return c.Internal.Request("GET", path, nil)
}

// Head :
func (c *Client) Head(path string) (Response, error, func()) {
	return c.Internal.Request("HEAD", path, nil)
}

// Post :
func (c *Client) Post(
	path, contentType string,
	body io.Reader,
) (Response, error, func()) {
	return c.Internal.Request("POST", path, body, func(req *http.Request) {
		req.Header.Set("Content-Type", contentType)
	})
}

// PostForm :
func (c *Client) PostForm(
	path string,
	data url.Values,
) (Response, error, func()) {
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
) (Response, error, func()) {
	return c.Post(path, "application/json", body)
}

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
