package tripperware

import (
	"net/http"
)

// RoundTripFunc :
type RoundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip :
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Ware :
type Ware func(http.RoundTripper) http.RoundTripper

// DecorateRoundTripper :
func (w Ware) DecorateRoundTripper(tripper http.RoundTripper) http.RoundTripper {
	return w(getTripper(tripper))
}

// RoundTripper :
func (w Ware) RoundTripper(req *http.Request) (*http.Response, error) {
	t := w.DecorateRoundTripper(http.DefaultTransport)
	return t.RoundTrip(req)
}

// DecorateClient :
func (w Ware) DecorateClient(client *http.Client, clone bool) *http.Client {
	c := getClient(client, clone)
	c.Transport = w.DecorateRoundTripper(c.Transport)
	return c
}

// Trippewares :
type List []Ware

// RoundTrip :
func (ws List) RoundTrip(req *http.Request) (*http.Response, error) {
	t := ws.DecorateRoundTripper(http.DefaultTransport)
	return t.RoundTrip(req)
}

// DecorateRoundTripper :
func (ws List) DecorateRoundTripper(tripper http.RoundTripper) http.RoundTripper {
	t := getTripper(tripper)
	for _, w := range ws {
		t = w(t)
	}
	return t
}

// DecorateClient :
func (ws List) DecorateClient(client *http.Client, clone bool) *http.Client {
	c := getClient(client, clone)
	c.Transport = ws.DecorateRoundTripper(c.Transport)
	return c
}

// Stack :
func Stack(wares ...Ware) List {
	return wares
}

func getTripper(tripper http.RoundTripper) http.RoundTripper {
	if tripper == nil {
		return http.DefaultTransport
	}
	return tripper
}

func getClient(client *http.Client, clone bool) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	if clone {
		c := *client
		client = &c
	}
	return client
}
