package testclient

import "io"

// Config :
type Config struct {
	BasePath string
	Body     io.Reader

	Decorator RoundTripperDecorator
}

// Copy :
func (c *Config) Copy() *Config {
	new := *c
	return &new
}
