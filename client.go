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
	return c.Internal.Request(method, path, body, func(req *http.Request) {
		for _, modify := range config.RequestModifiers {
			modify(req)
		}
	})
}

// DoFromRequest :
func (c *Client) DoFromRequest(
	req *http.Request,
) (Response, error, func()) {
	for _, modify := range c.Config.RequestModifiers {
		modify(req)
	}
	return c.Internal.Do(req)
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

// AddModifyRequest :
func AddModifyRequest(modify func(*http.Request)) func(*Config) {
	return func(c *Config) {
		c.RequestModifiers = append(c.RequestModifiers, modify)
	}
}
