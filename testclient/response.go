package testclient

import (
	"fmt"
	"net/http"
)

// Response :
type Response interface {
	Close()

	Response() *http.Response
	StatusCode() int

	Extractor
}

// Extractor :
type Extractor interface {
	ParseJSONData(val interface{}) error
	JSONData() interface{}

	Body() []byte
	LazyBodyString() fmt.Stringer
}
