package webtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/podhmo/go-webtest/testclient"
)

// Response :
type Response = testclient.Response

// Middleware :
type Middleware = func(
	t testing.TB,
	req *http.Request,
	inner func(*http.Request) (Response, error, func()),
) (Response, error, func())

// Internal :
type Internal interface {
	Do(req *http.Request) (Response, error, func())
	NewRequest(method string, path string, body io.Reader) (*http.Request, error)
}

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

// Bind :
func (c *Client) Bind(options ...func(*Config)) *Client {
	newClient := c.Copy()
	for _, opt := range options {
		opt(newClient.Config)
	}
	return newClient
}

// Do :
func (c *Client) Do(
	t testing.TB,
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
	return c.communicate(t, req, config)
}

// DoFromRequest :
func (c *Client) DoFromRequest(
	t testing.TB,
	req *http.Request,
	options ...func(*Config),
) (Response, error, func()) {
	config := c.Config.Copy()
	for _, opt := range options {
		opt(config)
	}
	return c.communicate(t, req, config)
}

// DoFromRequest :
func (c *Client) communicate(
	t testing.TB,
	req *http.Request,
	config *Config,
) (Response, error, func()) {
	for _, transform := range config.Transformers {
		transform(req)
	}

	doRequet := c.Internal.Do
	for i := range config.Middlewares {
		middleware := config.Middlewares[i]
		inner := doRequet
		doRequet = func(req *http.Request) (Response, error, func()) {
			return middleware(t, req, inner)
		}
	}
	return doRequet(req)
}

// NewClientFromTestServer :
func NewClientFromTestServer(ts *httptest.Server, options ...func(*Config)) *Client {
	c := NewConfig()
	for _, opt := range options {
		opt(c)
	}
	return &Client{
		Internal: &testclient.ServerClient{
			Server:    ts,
			BasePath:  c.BasePath,
			Transport: c.RoundTripper,
		},
		Config: c,
	}
}

// NewClientFromHandler :
func NewClientFromHandler(handler http.Handler, options ...func(*Config)) *Client {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}
	return &Client{
		Internal: &testclient.RecorderClient{
			Handler:  handler,
			BasePath: config.BasePath,
		},
		Config: config,
	}
}

// Config :
type Config struct {
	BasePath string
	Method   string

	Transformers []func(*http.Request) // request transformers
	Middlewares  []Middleware          // client middlewares
	RoundTripper http.RoundTripper

	body io.Reader // only once
}

// NewConfig :
func NewConfig() *Config {
	return &Config{
		Method: "GET",
	}
}

// Copy :
func (c *Config) Copy() *Config {
	return &Config{
		BasePath: c.BasePath,
		Transformers: append(
			make([]func(*http.Request), 0, len(c.Transformers)),
			c.Transformers...,
		),
		Middlewares: append(
			make([]Middleware, 0, len(c.Middlewares)),
			c.Middlewares...,
		),
	}
}

// WithMethod set method
func WithMethod(method string) func(*Config) {
	return func(c *Config) {
		c.Method = method
	}
}

// WithBasePath set base path
func WithBasePath(basePath string) func(*Config) {
	return func(c *Config) {
		c.BasePath = basePath
	}
}

// WithForm setup as send form-data request
func WithForm(data url.Values) func(*Config) {
	return func(c *Config) {
		if c.body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.body = strings.NewReader(data.Encode())
		c.Transformers = append(c.Transformers, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		})
	}
}

// WithJSON setup as json request
func WithJSON(body io.Reader) func(*Config) {
	return func(c *Config) {
		if c.body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.body = body
		c.Transformers = append(c.Transformers, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		})
	}
}

// WithTransformer adds request transformer
func WithTransformer(transform func(*http.Request)) func(*Config) {
	return func(c *Config) {
		c.Transformers = append(c.Transformers, transform)
	}
}

// WithRoundTripper adds request transformer
func WithRoundTripper(roundTripper http.RoundTripper) func(*Config) {
	return func(c *Config) {
		c.RoundTripper = roundTripper
	}
}
