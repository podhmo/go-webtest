package webtest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/podhmo/go-webtest/testclient"
	"github.com/podhmo/go-webtest/tripperware"
)

// Option :
type Option interface {
	Apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) Apply(c *Config) {
	f(c)
}

// Response :
type Response = testclient.Response

// Internal :
type Internal interface {
	Do(req *http.Request, clientConfig *testclient.Config) (Response, error)
	NewRequest(method string, path string, clientConfig *testclient.Config) (*http.Request, error)
}

// Client :
type Client struct {
	Internal Internal
}

// Get :
func (c *Client) Get(path string, options ...Option) (Response, error) {
	return c.Do("GET", path, options...)
}

// Post :
func (c *Client) Post(path string, options ...Option) (Response, error) {
	return c.Do("POST", path, options...)
}

// Put :
func (c *Client) Put(path string, options ...Option) (Response, error) {
	return c.Do("PUT", path, options...)
}

// Patch :
func (c *Client) Patch(path string, options ...Option) (Response, error) {
	return c.Do("PATCH", path, options...)
}

// Delete :
func (c *Client) Delete(path string, options ...Option) (Response, error) {
	return c.Do("DELETE", path, options...)
}

// Head :
func (c *Client) Head(path string, options ...Option) (Response, error) {
	return c.Do("HEAD", path, options...)
}

// Do :
func (c *Client) Do(
	method string,
	path string,
	options ...Option,
) (Response, error) {
	config := NewConfig()
	for _, opt := range options {
		opt.Apply(config)
	}

	req, err := c.Internal.NewRequest(method, path, config.ClientConfig)
	if err != nil {
		return nil, err
	}
	return c.communicate(req, config)
}

// DoFromRequest :
func (c *Client) DoFromRequest(
	req *http.Request,
	options ...Option,
) (Response, error) {
	config := NewConfig()
	for _, opt := range options {
		opt.Apply(config)
	}
	return c.communicate(req, config)
}

// DoFromRequest :
func (c *Client) communicate(
	req *http.Request,
	config *Config,
) (Response, error) {
	for _, modifyRequest := range config.ModifyRequests {
		modifyRequest(req)
	}

	doRequest := func(req *http.Request) (Response, error) {
		return c.Internal.Do(req, config.ClientConfig)
	}
	return doRequest(req)
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
		ClientConfig: c.ClientConfig.Copy(),
	}
}

// WithBasePath set base path
func WithBasePath(basePath string) Option {
	return optionFunc(func(c *Config) {
		c.ClientConfig.BasePath = basePath
	})
}

// WithQuery :
func WithQuery(query url.Values) Option {
	return optionFunc(func(c *Config) {
		c.ClientConfig.Query = query
	})
}

// WithForm setup as send form-data request
func WithForm(data url.Values) Option {
	return optionFunc(func(c *Config) {
		if c.ClientConfig.Body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.ClientConfig.Body = strings.NewReader(data.Encode())
		c.ModifyRequests = append(c.ModifyRequests, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		})
	})
}

// WithJSON setup as json request
func WithJSON(body io.Reader) Option {
	return optionFunc(func(c *Config) {
		if c.ClientConfig.Body != nil {
			panic("body is already set, enable to set body only once") // xxx
		}
		c.ClientConfig.Body = body
		c.ModifyRequests = append(c.ModifyRequests, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		})
	})
}

// WithModifyRequest adds request modifyRequest
func WithModifyRequest(modifyRequest func(*http.Request)) Option {
	return optionFunc(func(c *Config) {
		c.ModifyRequests = append(c.ModifyRequests, modifyRequest)
	})
}

// WithTripperware with client side middleware for roundTripper
func WithTripperware(wares ...tripperware.Ware) Option {
	return optionFunc(func(c *Config) {
		c.ClientConfig.Tripperwares = append(c.ClientConfig.Tripperwares, wares...)
	})
}
