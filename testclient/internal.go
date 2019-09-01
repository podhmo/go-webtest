package testclient

import (
	"log"
	"net/http"
	"os"
)

var (
	defaultDecorator      RoundTripperDecorator
	defaultInternalClient *http.Client
)

// TODO: logging interface

func init() {
	if os.Getenv("DEBUG") == "" {
		defaultInternalClient = http.DefaultClient
		return
	} else {
		copied := *http.DefaultClient
		defaultDecorator = NewDebugRoundTripper()
		copied.Transport = defaultDecorator.Decorate(copied.Transport)
		defaultInternalClient = &copied
		log.Println("builtin DebugRoundTripper is activated")
	}
}

func getInternalClient(client *http.Client, decorator RoundTripperDecorator) *http.Client {
	if client == nil {
		client = defaultInternalClient
	}

	if decorator != nil {
		// shallow copy
		copied := *client
		copied.Transport = getDecoratepedTransport(copied.Transport, decorator)
		return &copied
	}
	return client
}

func getDecoratepedTransport(original http.RoundTripper, decorator RoundTripperDecorator) http.RoundTripper {
	if decorator == nil {
		// if not debug, defaultDecorator is nil
		if defaultDecorator == nil {
			return original
		}
		decorator = defaultDecorator
	}
	return decorator.Decorate(original)
}
