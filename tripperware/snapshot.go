package tripperware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/podhmo/go-webtest/snapshot"
	"github.com/podhmo/go-webtest/tripperware/internal"
)

// GetExpectedDataFromSnapshot :
func GetExpectedDataFromSnapshot(
	t testing.TB,
	want *interface{},
	options ...func(sc *snapshot.Config),
) Ware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			e := newEmitter()
			if err := e.OnRequest(req); err != nil {
				return nil, err
			}
			res, err := next.RoundTrip(req)
			copied := internal.CopyResponse(res)
			err = e.OnResponse(res, err)

			// assign (side-effect!!), want is response data
			loaddata := snapshot.Take(t, e.Data, options...)
			*want = loaddata.(map[string]interface{})["response"].(map[string]interface{})["data"]
			return copied, err
		})
	}
}

// TakeSnapshot always takes a snapshot
func TakeSnapshot(
	t testing.TB,
	options ...func(sc *snapshot.Config),
) Ware {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			e := newEmitter()
			if err := e.OnRequest(req); err != nil {
				return nil, err
			}
			res, err := next.RoundTrip(req)
			copied := internal.CopyResponse(res)
			err = e.OnResponse(res, err)

			_ = snapshot.Take(t,
				e.Data,
				append([]func(*snapshot.Config){snapshot.WithForceUpdate()}, options...)...,
			)
			return copied, err
		})
	}
}

type emitter struct {
	Data map[string]map[string]interface{}
}

func newEmitter() *emitter {
	return &emitter{
		Data: map[string]map[string]interface{}{
			"request":  map[string]interface{}{},
			"response": map[string]interface{}{},
		},
	}
}

func (e *emitter) OnRequest(req *http.Request) error {
	e.Data["request"]["method"] = req.Method
	e.Data["request"]["path"] = req.URL.Path + req.URL.RawQuery

	hasBody := req.Body != nil && req.Body != http.NoBody
	if hasBody {
		var b bytes.Buffer
		var payload interface{}
		r := io.TeeReader(req.Body, &b)
		decoder := json.NewDecoder(r) // TODO: see Content-type
		if err := decoder.Decode(&payload); err != nil {
			return err
		}
		e.Data["request"]["body"] = payload
		req.Body = ioutil.NopCloser(&b)
	}
	return nil
}

func (e *emitter) OnResponse(res *http.Response, err error) error {
	var data interface{}
	if err == nil {
		decoder := json.NewDecoder(res.Body) // TODO: see Content-type
		defer res.Body.Close()
		err = decoder.Decode(&data)
	}
	if err != nil {
		e.Data["error"] = map[string]interface{}{
			"message": err.Error(),
			// "verbose": fmt.Sprintf("%+v", err),
		}
	}
	if res != nil {
		e.Data["response"]["statusCode"] = res.StatusCode
		if data != nil {
			e.Data["response"]["data"] = data
		}
	}
	return err
}
