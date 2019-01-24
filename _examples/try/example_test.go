package examples

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/podhmo/go-webtest/try"
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
		try.JSONRequest(
			t,
			mux,
			"GET",
			"/",
			http.StatusOK,
			try.WithAssertJSONResponse(
				`{"message": "hello world"}`,
			))
	})

	t.Run("response mismatch", func(t *testing.T) {
		mux := http.HandlerFunc(Handler)
		try.Request(
			t,
			mux,
			"GET",
			"/",
			http.StatusOK,
			try.WithJSONBody(
				`{"name": "WORLD"}`,
			),
			try.WithAssertJSONResponse(
				`{"message": "hello world"}`,
			))
	})
}
