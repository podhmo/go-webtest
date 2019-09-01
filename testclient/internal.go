package testclient

import (
	"log"
	"net/http"
	"os"
)

var (
	defaultTransport      RoundTripperDecorator
	defaultInternalClient *http.Client
)

// TODO: logging interface

func init() {
	if os.Getenv("DEBUG") == "" {
		defaultInternalClient = http.DefaultClient
		return
	} else {
		copied := *http.DefaultClient
		defaultTransport = &DebugRoundTripper{}
		copied.Transport = defaultTransport.Decorate(copied.Transport)
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
		copied.Transport = getDecoratepedTransport(copied.Transport, transport)
		return &copied
	}
	return client
}

func getDecoratepedTransport(original, transport http.RoundTripper) http.RoundTripper {
	if transport == nil {
		transport = defaultTransport // if not debug, defaultTransport is nil
	}
	if transport == nil {
		return original
	}
	if t, ok := transport.(RoundTripperDecorator); ok {
		return t.Decorate(original)
	}
	log.Printf("client.Transport is already set, config.Transport[%T] is ignored", transport)
	return original
}
