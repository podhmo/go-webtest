package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/podhmo/go-webtest/internal"
)

// NewAdapter :
func NewAdapter(get func() *http.Response) *Adapter {
	return &Adapter{
		GetResponse: get,
	}
}

// Adapter :
type Adapter struct {
	GetResponse func() *http.Response

	bytes []byte // not bet
	bOnce sync.Once

	m         sync.Mutex
	teardowns []func() error
}

// Response :
func (res *Adapter) Response() *http.Response {
	return res.GetResponse()
}

// Close :
func (res *Adapter) Close() {
	res.m.Lock()
	defer res.m.Unlock()
	for _, teardown := range res.teardowns {
		if err := teardown(); err != nil {
			panic(err)
		}
	}
	res.teardowns = nil
}

func (res *Adapter) AddTeardown(fn func() error) {
	res.m.Lock()
	defer res.m.Unlock()
	res.teardowns = append(res.teardowns, fn)
}

// Buffer : (TODO: rename)
func (res *Adapter) Buffer() *bytes.Buffer {
	res.bOnce.Do(func() {
		var b bytes.Buffer
		if _, err := io.Copy(&b, res.Response().Body); err != nil {
			panic(err) // xxx
		}
		res.bytes = b.Bytes()
	})
	return bytes.NewBuffer(res.bytes)
}

// StatusCode :
func (res *Adapter) StatusCode() int {
	return res.Response().StatusCode
}

// ParseJSONData :
func (res *Adapter) ParseJSONData(val interface{}) error {
	decoder := json.NewDecoder(res.Buffer()) // TODO: decoder interface
	return decoder.Decode(val)
}

// JSONData :
func (res *Adapter) JSONData() interface{} {
	var val interface{}
	if err := res.ParseJSONData(&val); err != nil {
		panic(err) // xxx:
	}
	return val
}

// Body :
func (res *Adapter) Body() []byte {
	return res.Buffer().Bytes()
}

// LazyBodyString :
func (res *Adapter) LazyBodyString() fmt.Stringer {
	return internal.NewLazyString(
		func() string {
			return res.Buffer().String()
		},
	)
}
