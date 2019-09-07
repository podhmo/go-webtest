package testclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/podhmo/go-webtest/testclient/internal"
)

// Response :
type Response interface {
	io.Closer

	Code() int
	Header() http.Header
	Request() *http.Request

	Raw() *http.Response

// 	Extractor
// }

// // Extractor :
// type Extractor interface {
	ParseJSONData(val interface{}) error
	JSONData() interface{}

	Body() []byte
	Text() string
	LazyText() fmt.Stringer
}

// NewResponseAdapter :
func NewResponseAdapter(get func() *http.Response) *ResponseAdapter {
	return &ResponseAdapter{
		GetResponse: get,
	}
}

// ResponseAdapter :
type ResponseAdapter struct {
	GetResponse func() *http.Response

	bytes []byte // not bet
	bOnce sync.Once

	m         sync.Mutex
	teardowns []func() error
}

// Raw :
func (res *ResponseAdapter) Raw() *http.Response {
	return res.GetResponse()
}

// Close :
func (res *ResponseAdapter) Close() error {
	res.m.Lock()
	defer res.m.Unlock()
	for _, teardown := range res.teardowns {
		if err := teardown(); err != nil {
			return err
		}
	}
	res.teardowns = nil
	return nil
}

func (res *ResponseAdapter) AddTeardown(fn func() error) {
	res.m.Lock()
	defer res.m.Unlock()
	res.teardowns = append(res.teardowns, fn)
}

// Buffer : (TODO: rename)
func (res *ResponseAdapter) Buffer() *bytes.Buffer {
	res.bOnce.Do(func() {
		var b bytes.Buffer
		if _, err := io.Copy(&b, res.Raw().Body); err != nil {
			panic(err) // xxx
		}
		res.bytes = b.Bytes()
	})
	return bytes.NewBuffer(res.bytes)
}

// Code :
func (res *ResponseAdapter) Code() int {
	return res.Raw().StatusCode
}

// Header :
func (res *ResponseAdapter) Header() http.Header {
	return res.Raw().Header
}

// Request :
func (res *ResponseAdapter) Request() *http.Request {
	return res.Raw().Request
}

// ParseJSONData :
func (res *ResponseAdapter) ParseJSONData(val interface{}) error {
	decoder := json.NewDecoder(res.Buffer()) // TODO: decoder interface
	return decoder.Decode(val)
}

// JSONData :
func (res *ResponseAdapter) JSONData() interface{} {
	var val interface{}
	if err := res.ParseJSONData(&val); err != nil {
		panic(err) // xxx:
	}
	return val
}

// Body :
func (res *ResponseAdapter) Body() []byte {
	return res.Buffer().Bytes()
}

// Text :
func (res *ResponseAdapter) Text() string {
	return res.Buffer().String()
}

// LazyText :
func (res *ResponseAdapter) LazyText() fmt.Stringer {
	return internal.NewLazyString(
		func() string {
			return res.Buffer().String()
		},
	)
}
