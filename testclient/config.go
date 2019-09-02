package testclient

import (
	"io"
	"net/url"

	"github.com/podhmo/go-webtest/tripperware"
)

// Config :
type Config struct {
	BasePath string
	Body     io.Reader
	Query    url.Values

	Tripperwares tripperware.List
}

// Copy :
func (c *Config) Copy() *Config {
	new := *c
	new.Tripperwares = append(
		make([]tripperware.Ware, 0, len(c.Tripperwares)),
		c.Tripperwares...,
	)
	return &new
}
