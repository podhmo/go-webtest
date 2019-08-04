package middlewares

import (
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/snapshot"
)

// SnapshotTesting :
func SnapshotTesting(options ...func(sc *snapshot.Config)) func(*webtest.Config) {
	return func(c *webtest.Config) {
		c.Middlewares = append(c.Middlewares, NewMiddleware(
			func(
				t testing.TB,
				res Response,
				req *http.Request,
			) error {
				// TODO: get stored data, in testcode
				// TODO: following .har structure?

				storedata := map[string]interface{}{
					"request": map[string]interface{}{
						"method": req.Method,
						"path":   req.URL.Path + req.URL.RawQuery,
					},
					"response": map[string]interface{}{
						"statusCode": res.StatusCode(),
						"data":       res.JSONData(),
					},
				}
				_ = snapshot.Take(t, storedata, options...)
				return nil
			}))
	}
}
