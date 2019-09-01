package testclient

import (
	"log"
	"net/http"
	"os"
)

var (
	defaultTransport      roundTripperWrapper
	defaultInternalClient *http.Client
)

// TODO: logging interface

type roundTripperWrapper interface {
	http.RoundTripper
	Wrap(inner http.RoundTripper) http.RoundTripper
}

func init() {
	if os.Getenv("DEBUG") == "" {
		defaultInternalClient = http.DefaultClient
		return
	} else {
		copied := *http.DefaultClient
		defaultTransport = &DebugRoundTripper{}
		copied.Transport = defaultTransport.Wrap(copied.Transport)
		defaultInternalClient = &copied
		log.Println("builtin DebugRoundTripper is activated")
	}
}

func getInternalClientWithTransport(client *http.Client, transport http.RoundTripper) *http.Client {
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
		copied.Transport = getWrappedTransport(copied.Transport, transport)
		return &copied
	}
	return client
}

func getWrappedTransport(original, transport http.RoundTripper) http.RoundTripper {
	if transport == nil {
		transport = defaultTransport // if not debug, defaultTransport is nil
	}
	if transport == nil {
		return original
	}
	if t, ok := transport.(roundTripperWrapper); ok {
		return t.Wrap(original)
	}
	log.Printf("client.Transport is already set, config.Transport[%T] is ignored", transport)
	return original
}
