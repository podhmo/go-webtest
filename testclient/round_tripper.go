package testclient

import (
	"log"
	"net/http"
)

// RoundTripperDecorateFunc :
type RoundTripperDecorateFunc func(http.RoundTripper, *http.Request) (*http.Response, error)

// Decorate :
func (f RoundTripperDecorateFunc) Decorate(inner http.RoundTripper) RoundTripperDecorator {
	return FuncRoundTripper{Fn: f}.Decorate(inner)
}

// RoundTripperDecorator :
type RoundTripperDecorator interface {
	http.RoundTripper
	Decorate(inner http.RoundTripper) RoundTripperDecorator
}

// FuncRoundTripper :
type FuncRoundTripper struct {
	Fn           RoundTripperDecorateFunc
	RoundTripper http.RoundTripper
}

// Decorate :
func (w FuncRoundTripper) Decorate(inner http.RoundTripper) RoundTripperDecorator {
	if inner, ok := inner.(RoundTripperDecorator); ok {
		return FuncRoundTripper{
			Fn: func(tripper http.RoundTripper, req *http.Request) (*http.Response, error) {
				return w.Fn(inner, req)
			},
		}
	}
	if w.RoundTripper != nil {
		log.Printf("!! %T.RoundTripper is not nil, overwrite original one", w)
	}
	return FuncRoundTripper{Fn: w.Fn, RoundTripper: inner}
}

// RoundTrip :
func (w FuncRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	tripper := w.RoundTripper
	if tripper == nil {
		tripper = http.DefaultTransport // xxx:
	}
	return w.Fn(tripper, req)
}