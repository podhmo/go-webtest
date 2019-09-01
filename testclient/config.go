package testclient

import (
	"io"
	"net/url"
)

// Config :
type Config struct {
	BasePath string
	Body     io.Reader
	Query    url.Values

	Decorator RoundTripperDecorator
}

// Copy :
func (c *Config) Copy() *Config {
	new := *c
	return &new
}
