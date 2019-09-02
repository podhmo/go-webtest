package hook

import (
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/snapshot"
)

// GetExpectedDataFromSnapshot :
func GetExpectedDataFromSnapshot(
	t testing.TB,
	want *interface{},
	options ...func(sc *snapshot.Config),
) webtest.Option {
	return func(c *webtest.Config) {
		c.Hooks = append(c.Hooks, NewHook(
			func(
				res Response,
				req *http.Request,
			) error {
				storedata := createSnapshotData(res, req)

				// assign (side-effect!!), want is response data
				loaddata := snapshot.Take(t, storedata, options...)
				*want = loaddata.(map[string]interface{})["response"].(map[string]interface{})["data"]
				return nil
			}))
	}
}

// TakeSnapshot always takes a snapshot
func TakeSnapshot(
	t testing.TB,
	options ...func(sc *snapshot.Config),
) webtest.Option {
	return func(c *webtest.Config) {
		c.Hooks = append(c.Hooks, NewHook(
			func(
				res Response,
				req *http.Request,
			) error {
				storedata := createSnapshotData(res, req)
				_ = snapshot.Take(t,
					storedata,
					append([]func(*snapshot.Config){snapshot.WithForceUpdate()}, options...)...,
				)
				return nil
			}))
	}
}

func createSnapshotData(res Response, req *http.Request) interface{} {
	// TODO: following .har structure? or openapi spec structure?
	return map[string]interface{}{
		"request": map[string]interface{}{
			"method": req.Method,
			"path":   req.URL.Path + req.URL.RawQuery,
		},
		"response": map[string]interface{}{
			"statusCode": res.Code(),
			"data":       res.JSONData(),
		},
	}
}
