package testclient

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

type roundTripperWrapper struct {
	http.RoundTripper
	Wrap func(inner http.RoundTripper) http.RoundTripper
}

// TODO: logging interface

// ServerClient :
type ServerClient struct {
	Client    *http.Client
	Transport http.RoundTripper

	Server   *httptest.Server
	BasePath string // need?
}

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

func (c *ServerClient) client() *http.Client {
	client := c.Client
	if client == nil {
		client = defaultInternalClient
	}

	if c.Transport != nil {
		// shallow copy
		copied := *client
		if copied.Transport == nil {
			copied.Transport = c.Transport
			return &copied
		}
		switch t := c.Transport.(type) {
		case roundTripperWrapper:
			copied.Transport = t.Wrap(copied.Transport)
			return &copied
		default:
			log.Printf("client.Transport is already set, config.Transport[%T] is ignored", c.Transport)
		}
	}
	return client
}

// Do :
func (c *ServerClient) Do(req *http.Request) (Response, error, func()) {
	client := c.client()

	var adapter *ResponseAdapter
	var raw *http.Response
	var once sync.Once

	raw, err := client.Do(req)
	if err != nil {
		return nil, err, nil
	}

	adapter = NewResponseAdapter(
		func() *http.Response {
			once.Do(func() {
				adapter.AddTeardown(raw.Body.Close)
			})
			return raw
		},
	)
	return adapter, err, adapter.Close
}

// NewRequest :
func (c *ServerClient) NewRequest(
	method string,
	path string,
	body io.Reader,
) (*http.Request, error) {
	url := internal.URLJoin(c.Server.URL, internal.URLJoin(c.BasePath, path))
	return http.NewRequest(method, url, body)
}
