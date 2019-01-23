package examples

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "/application/json")
	decoder := json.NewDecoder(r.Body)

	params := map[string]string{}
	if err := decoder.Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "invalid value"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "hello %[1]s", "name": "%[1]s"}`, params["name"])
}

func Test(t *testing.T) {
	t.Run("status mismatch", func(t *testing.T) {
		mux := http.HandlerFunc(Handler)
		webtest.TryJSONRequest(
			t,
			mux,
			"GET",
			"/",
			http.StatusOK,
			webtest.WithAssertJSONResponse(
				`{"message": "hello world"}`,
			))
	})

	t.Run("response mismatch", func(t *testing.T) {
		mux := http.HandlerFunc(Handler)
		webtest.TryRequest(
			t,
			mux,
			"GET",
			"/",
			http.StatusOK,
			webtest.WithJSONBody(
				`{"name": "WORLD"}`,
			),
			webtest.WithAssertJSONResponse(
				`{"message": "hello world"}`,
			))
	})
}
