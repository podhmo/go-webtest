package webtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/podhmo/go-webtest/client"
	"github.com/podhmo/go-webtest/client/response"
)

// Client :
type Client interface {
	Do(req *http.Request) (Response, error, func())
	Get(path string) (Response, error, func())
	Head(path string) (Response, error, func())
	Post(path, contentType string, body io.Reader) (Response, error, func())
	PostJSON(path string, body io.Reader) (Response, error, func())
	PostForm(path string, data url.Values) (Response, error, func())
	// todo: setting by functional options?
}

// Response :
type Response = response.Response

// NewClientFromTestServer :
func NewClientFromTestServer(ts *httptest.Server, options ...func(*Config)) Client {
	c := &Config{}
	for _, opt := range options {
		opt(c)
	}
	return &client.Adapter{
		Internal: &client.HTTPTestServerClient{
			Server:   ts,
			BasePath: c.BasePath,
		},
	}
}

// NewClientFromHandler :
func NewClientFromHandler(handlerFunc http.HandlerFunc, options ...func(*Config)) Client {
	c := &Config{}
	for _, opt := range options {
		opt(c)
	}
	return &client.Adapter{
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
