package middlewares

import (
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/snapshot"
)

// SnapshotTesting :
func SnapshotTesting(want *interface{}, options ...func(sc *snapshot.Config)) func(*webtest.Config) {
	return func(c *webtest.Config) {
		c.Middlewares = append(c.Middlewares, NewMiddleware(
			func(
				t testing.TB,
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
func TakeSnapshot(options ...func(sc *snapshot.Config)) func(*webtest.Config) {
	return func(c *webtest.Config) {
		c.Middlewares = append(c.Middlewares, NewMiddleware(
			func(
				t testing.TB,
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
			"statusCode": res.StatusCode(),
			"data":       res.JSONData(),
		},
	}
}
