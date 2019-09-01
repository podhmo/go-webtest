package testclient

import (
	"log"
	"net/http"
	"os"
)

var (
	defaultInternalClient *http.Client
)

func init() {
	if os.Getenv("DEBUG") == "" {
		defaultInternalClient = http.DefaultClient
		return
	} else {
		copied := *http.DefaultClient
		copied.Transport = &DebugRoundTripper{}
		defaultInternalClient = &copied
		log.Println("builtin DebugRoundTripper is activated")
	}
}

// GetInternalClientWith :
func GetInternalClientWith(client *http.Client, transport http.RoundTripper) *http.Client {
	if client == nil {
		client = defaultInternalClient
	}

	if transport != nil {
		// shallow copy
		copied := *client
		if copied.Transport == nil {
			copied.Transport = transport
			return &copied
		}
		switch t := transport.(type) {
		case roundTripperWrapper:
			copied.Transport = t.Wrap(copied.Transport)
			return &copied
		default:
			log.Printf("client.Transport is already set, config.Transport[%T] is ignored", transport)
		}
	}
	return client
}
