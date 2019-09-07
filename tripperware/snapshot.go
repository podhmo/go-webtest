package tripperware

import (
	"encoding/json"
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
			res, err := next.RoundTrip(req)
			copied := internal.CopyResponse(res)
			storedata, err := createSnapshotData(res, req, err)

			// assign (side-effect!!), want is response data
			loaddata := snapshot.Take(t, storedata, options...)
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
			res, err := next.RoundTrip(req)
			copied := internal.CopyResponse(res)
			storedata, err := createSnapshotData(res, req, err)
			_ = snapshot.Take(t,
				storedata,
				append([]func(*snapshot.Config){snapshot.WithForceUpdate()}, options...)...,
			)
			return copied, err
		})
	}
}

func createSnapshotData(res *http.Response, req *http.Request, err error) (map[string]map[string]interface{}, error) {
	var data interface{}
	if err == nil {
		decoder := json.NewDecoder(res.Body) // TODO: see Content-type
		defer res.Body.Close()
		err = decoder.Decode(&data)
	}
	storedata := map[string]map[string]interface{}{
		"response": map[string]interface{}{},
	}
	if err != nil {
		storedata["error"] = map[string]interface{}{"message": err.Error()}
	}
	if req != nil {
		storedata["request"] = map[string]interface{}{
			"method": req.Method,
			"path":   req.URL.Path + req.URL.RawQuery,
		}
	}
	if res != nil {
		storedata["response"]["statusCode"] = res.StatusCode
		if data != nil {
			storedata["response"]["data"] = data
		}
	}
	return storedata, err
}
