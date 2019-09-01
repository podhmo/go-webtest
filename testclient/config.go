package testclient

// Config :
type Config struct {
	Decorator RoundTripperDecorator
}

// Copy :
func (c *Config) Copy() *Config {
	new := *c
	return &new
}
