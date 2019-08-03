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
	DoFromRequest(req *http.Request) (Response, error, func())
	NewRequest(method string, path string, body io.Reader) (*http.Request, error)
}

// todo: rename

// Middleware :
type Middleware = func(
	req *http.Request,
	inner func(*http.Request) (Response, error, func()),
) (Response, error, func())

// Client :
type Client struct {
	Internal Internal
	Config   *Config
}

// Copy :
func (c *Client) Copy() *Client {
	return &Client{
		Internal: c.Internal,
		Config:   c.Config.Copy(),
	}
}

// Do :
func (c *Client) Do(
	path string,
	options ...func(*Config),
) (Response, error, func()) {
	config := c.Config.Copy()
	for _, opt := range options {
		opt(config)
	}
	method := config.Method
	body := config.body
	req, err := c.Internal.NewRequest(method, path, body)
	if err != nil {
		return nil, err, nil
	}
	return c.do(req, config)
}

// DoFromRequest :
func (c *Client) DoFromRequest(
	req *http.Request,
	options ...func(*Config),
) (Response, error, func()) {
	config := c.Config.Copy()
	for _, opt := range options {
		opt(config)
	}
	return c.do(req, config)
}

// DoFromRequest :
func (c *Client) do(
	req *http.Request,
	config *Config,
) (Response, error, func()) {
	for _, modify := range config.RequestModifiers {
		modify(req)
	}

	doRequet := c.Internal.DoFromRequest
	for i := range config.Middlewares {
		middleware := config.Middlewares[i]
		inner := doRequet
		doRequet = func(req *http.Request) (Response, error, func()) {
			return middleware(req, inner)
		}
	}
	return doRequet(req)
}

// NewClientFromTestServer :
func NewClientFromTestServer(ts *httptest.Server, options ...func(*Config)) *Client {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}
	return &Client{
		Internal: &client.HTTPTestServerClient{
			Server:   ts,
			BasePath: config.BasePath,
		},
		Config: config,
	}
}

// NewClientFromHandler :
func NewClientFromHandler(handlerFunc http.HandlerFunc, options ...func(*Config)) *Client {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}
	return &Client{
		Internal: &client.HTTPTestResponseRecorderClient{
			HandlerFunc: handlerFunc,
			BasePath:    config.BasePath,
		},
		Config: config,
	}
}

// Config :
type Config struct {
	BasePath string

	Method           string
	RequestModifiers []func(*http.Request)
	Middlewares      []Middleware // todo: rename

	body io.Reader // only once
}

// NewConfig :
func NewConfig() *Config {
	return &Config{
		Method: "GET",
	}
}

// Config :
func (c *Config) Copy() *Config {
	return &Config{
		BasePath: c.BasePath,
		RequestModifiers: append(
			make([]func(*http.Request), 0, len(c.RequestModifiers)),
			c.RequestModifiers...,
		),
		Middlewares: append(
			make([]Middleware, 0, len(c.Middlewares)),
			c.Middlewares...,
		),
	}
}

// WithBasePath :
func WithBasePath(basePath string) func(*Config) {
	return func(c *Config) {
		c.BasePath = basePath
	}
}

// WithForm :
func WithForm(data url.Values) func(*Config) {
	return func(c *Config) {
		if c.body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.body = strings.NewReader(data.Encode())
		c.RequestModifiers = append(c.RequestModifiers, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		})
	}
}

// WithJSON :
func WithJSON(body io.Reader) func(*Config) {
	return func(c *Config) {
		if c.body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.body = body
		c.RequestModifiers = append(c.RequestModifiers, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		})
	}
}

// WithModifyRequest :
func WithModifyRequest(modify func(*http.Request)) func(*Config) {
	return func(c *Config) {
		c.RequestModifiers = append(c.RequestModifiers, modify)
	}
}

// WithMiddleware :
func WithMiddleware(middleware Middleware) func(*Config) {
	return func(c *Config) {
		c.Middlewares = append(c.Middlewares, middleware)
	}
}
