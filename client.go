package webtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/podhmo/go-webtest/testclient"
)

// Option :
type Option func(*Config)

// Response :
type Response = testclient.Response

// Hook :
type Hook = func(
	req *http.Request,
	inner func(*http.Request) (Response, error, func()),
) (Response, error, func())

// Internal :
type Internal interface {
	Do(req *http.Request, clientConfig *testclient.Config) (Response, error, func())
	NewRequest(method string, path string, clientConfig *testclient.Config) (*http.Request, error)
}

// Client :
type Client struct {
	Internal Internal
}

// GET :
func (c *Client) GET(path string, options ...Option) (Response, error, func()) {
	return c.Do("GET", path, options...)
}

// POST :
func (c *Client) POST(path string, options ...Option) (Response, error, func()) {
	return c.Do("POST", path, options...)
}

// PUT :
func (c *Client) PUT(path string, options ...Option) (Response, error, func()) {
	return c.Do("PUT", path, options...)
}

// PATCH :
func (c *Client) PATCH(path string, options ...Option) (Response, error, func()) {
	return c.Do("PATCH", path, options...)
}

// DELETE :
func (c *Client) DELETE(path string, options ...Option) (Response, error, func()) {
	return c.Do("DELETE", path, options...)
}

// HEAD :
func (c *Client) HEAD(path string, options ...Option) (Response, error, func()) {
	return c.Do("HEAD", path, options...)
}

// Do :
func (c *Client) Do(
	method string,
	path string,
	options ...Option,
) (Response, error, func()) {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}

	req, err := c.Internal.NewRequest(method, path, config.ClientConfig)
	if err != nil {
		return nil, err, nil
	}
	return c.communicate(req, config)
}

// DoFromRequest :
func (c *Client) DoFromRequest(
	req *http.Request,
	options ...Option,
) (Response, error, func()) {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}
	return c.communicate(req, config)
}

// DoFromRequest :
func (c *Client) communicate(
	req *http.Request,
	config *Config,
) (Response, error, func()) {
	for _, modifyRequest := range config.ModifyRequests {
		modifyRequest(req)
	}

	doRequet := func(req *http.Request) (Response, error, func()) {
		return c.Internal.Do(req, config.ClientConfig)
	}
	for i := range config.Hooks {
		hook := config.Hooks[i]
		inner := doRequet
		doRequet = func(req *http.Request) (Response, error, func()) {
			return hook(req, inner)
		}
	}
	return doRequet(req)
}

// NewClientFromTestServer :
func NewClientFromTestServer(ts *httptest.Server) *Client {
	return &Client{
		Internal: &testclient.RealClient{
			URL: ts.URL,
		},
	}
}

// NewClientFromHandler :
func NewClientFromHandler(handler http.Handler) *Client {
	return &Client{
		Internal: &testclient.FakeClient{
			Handler: handler,
		},
	}
}

// NewClientFromURL :
func NewClientFromURL(url string) *Client {
	return &Client{
		Internal: &testclient.RealClient{
			URL: url,
		},
	}
}

// Config :
type Config struct {
	BasePath     string
	Method       string
	ClientConfig *testclient.Config

	ModifyRequests []func(*http.Request) // request modifyRequests
	Hooks          []Hook                // client hooks
}

// NewConfig :
func NewConfig() *Config {
	return &Config{
		ClientConfig: &testclient.Config{},
	}
}

// Copy :
func (c *Config) Copy() *Config {
	return &Config{
		BasePath: c.BasePath,
		ModifyRequests: append(
			make([]func(*http.Request), 0, len(c.ModifyRequests)),
			c.ModifyRequests...,
		),
		Hooks: append(
			make([]Hook, 0, len(c.Hooks)),
			c.Hooks...,
		),
		ClientConfig: c.ClientConfig.Copy(),
	}
}

// WithBasePath set base path
func WithBasePath(basePath string) Option {
	return func(c *Config) {
		c.ClientConfig.BasePath = basePath
	}
}

// WithQuery :
func WithQuery(query url.Values) Option {
	return func(c *Config) {
		c.ClientConfig.Query = query
	}
}

// WithForm setup as send form-data request
func WithForm(data url.Values) Option {
	return func(c *Config) {
		if c.ClientConfig.Body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.ClientConfig.Body = strings.NewReader(data.Encode())
		c.ModifyRequests = append(c.ModifyRequests, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		})
	}
}

// WithJSON setup as json request
func WithJSON(body io.Reader) Option {
	return func(c *Config) {
		if c.ClientConfig.Body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.ClientConfig.Body = body
		c.ModifyRequests = append(c.ModifyRequests, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		})
	}
}

// WithModifyRequest adds request modifyRequest
func WithModifyRequest(modifyRequest func(*http.Request)) Option {
	return func(c *Config) {
		c.ModifyRequests = append(c.ModifyRequests, modifyRequest)
	}
}

// WithRoundTripperDecorator with client side middleware for roundTripper
func WithRoundTripperDecorator(decorator testclient.RoundTripperDecorator) Option {
	return func(c *Config) {
		c.ClientConfig.Decorator = decorator
	}
}
