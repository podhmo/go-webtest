package webtest

import (
	"net/http"
	"net/http/httptest"

	"github.com/podhmo/go-webtest/client"
	"github.com/podhmo/go-webtest/client/response"
)

// Client :
type Client interface {
	Do(req *http.Request) (Response, error, func())
	Get(path string) (Response, error, func())
}

// Response :
type Response = response.Response

// NewClientFromTestServer :
func NewClientFromTestServer(ts *httptest.Server, options ...func(*Config)) Client {
	c := &Config{}
	for _, opt := range options {
		opt(c)
	}
	return &client.HTTPTestServerClient{
		Server:   ts,
		BasePath: c.BasePath,
	}
}

// NewClientFromHandler :
func NewClientFromHandler(handlerFunc http.HandlerFunc, options ...func(*Config)) Client {
	c := &Config{}
	for _, opt := range options {
		opt(c)
	}
	return &client.HTTPTestResponseRecorderClient{
		HandlerFunc: handlerFunc,
		BasePath:    c.BasePath,
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
