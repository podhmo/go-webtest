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
	Do(req *http.Request, clientConfig *testclient.Config) (Response, error, func())
	NewRequest(method string, path string, clientConfig *testclient.Config) (*http.Request, error)
}

// Client :
type Client struct {
	Internal Internal
}

// GET :
func (c *Client) GET(t testing.TB, path string, options ...func(*Config)) (Response, error, func()) {
	return c.Do(t, "GET", path, options...)
}

// POST :
func (c *Client) POST(t testing.TB, path string, options ...func(*Config)) (Response, error, func()) {
	return c.Do(t, "POST", path, options...)
}

// PUT :
func (c *Client) PUT(t testing.TB, path string, options ...func(*Config)) (Response, error, func()) {
	return c.Do(t, "PUT", path, options...)
}

// PATCH :
func (c *Client) PATCH(t testing.TB, path string, options ...func(*Config)) (Response, error, func()) {
	return c.Do(t, "PATCH", path, options...)
}

// DELETE :
func (c *Client) DELETE(t testing.TB, path string, options ...func(*Config)) (Response, error, func()) {
	return c.Do(t, "DELETE", path, options...)
}

// HEAD :
func (c *Client) HEAD(t testing.TB, path string, options ...func(*Config)) (Response, error, func()) {
	return c.Do(t, "HEAD", path, options...)
}

// Do :
func (c *Client) Do(
	t testing.TB,
	method string,
	path string,
	options ...func(*Config),
) (Response, error, func()) {
	config := NewConfig()
	for _, opt := range options {
		opt(config)
	}

	req, err := c.Internal.NewRequest(method, path, config.ClientConfig)
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
	config := NewConfig()
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
	for _, modifyRequest := range config.ModifyRequests {
		modifyRequest(req)
	}

	doRequet := func(req *http.Request) (Response, error, func()) {
		return c.Internal.Do(req, config.ClientConfig)
	}
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
	Middlewares    []Middleware          // client middlewares
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
		Middlewares: append(
			make([]Middleware, 0, len(c.Middlewares)),
			c.Middlewares...,
		),
		ClientConfig: c.ClientConfig.Copy(),
	}
}

// WithBasePath set base path
func WithBasePath(basePath string) func(*Config) {
	return func(c *Config) {
		c.ClientConfig.BasePath = basePath
	}
}

// WithQuery :
func WithQuery(query url.Values) func(*Config) {
	return func(c *Config) {
		c.ClientConfig.Query = query
	}
}

// WithForm setup as send form-data request
func WithForm(data url.Values) func(*Config) {
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
func WithJSON(body io.Reader) func(*Config) {
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
func WithModifyRequest(modifyRequest func(*http.Request)) func(*Config) {
	return func(c *Config) {
		c.ModifyRequests = append(c.ModifyRequests, modifyRequest)
	}
}

// WithRoundTripperDecorator with client side middleware for roundTripper
func WithRoundTripperDecorator(decorator testclient.RoundTripperDecorator) func(*Config) {
	return func(c *Config) {
		c.ClientConfig.Decorator = decorator
	}
}
