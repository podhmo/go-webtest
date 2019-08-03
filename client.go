package webtest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Client :
type Client interface {
	Do(req *http.Request) (Response, error, func())
	Get(path string) (Response, error, func())
}

// Response :
type Response interface {
	Close()

	Response() *http.Response
	StatusCode() int

	Extractor
}

// Extractor :
type Extractor interface {
	ParseData(val interface{}) error
	Data() interface{}

	Body() []byte
	LazyBodyString() fmt.Stringer
}

// NewClientFromTestServer :
func NewClientFromTestServer(ts *httptest.Server, options ...func(*Config)) Client {
	c := &Config{}
	for _, opt := range options {
		opt(c)
	}
	return &HTTPTestServerClient{
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
	return &HTTPTestResponseRecorderClient{
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
