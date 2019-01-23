package webtest_test

import (
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
)

func TestTryRequest(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		handler := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "ok"}`))
			},
		)
		webtest.TryJSONRequest(
			t,
			handler,
			"GET",
			"/",
			http.StatusOK,
			webtest.WithAssertJSONResponse(`{"message": "ok"}`),
		)
	})
}
