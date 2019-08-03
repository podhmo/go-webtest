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

// Requset :
func (c *Client) Request(
	path string,
	options ...func(*Config),
) (Response, error, func()) {
	config := c.Config.Copy()
	for _, opt := range options {
		opt(config)
	}
	method := config.Method
	body := config.body
	return c.Internal.Request(method, path, body)
}

// Do :
func (c *Client) Do(
	req *http.Request,
) (Response, error, func()) {
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

	Method          string
	RequstModifiers []func(*http.Request)

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
		RequstModifiers: append(
			make([]func(*http.Request), 0, len(c.RequstModifiers)),
			c.RequstModifiers...,
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
		c.RequstModifiers = append(c.RequstModifiers, func(req *http.Request) {
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
		c.RequstModifiers = append(c.RequstModifiers, func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		})
	}
}

// AddModifyRequest :
func AddModifyRequest(modify func(*http.Request)) func(*Config) {
	return func(c *Config) {
		c.RequstModifiers = append(c.RequstModifiers, modify)
	}
}
