package try_test

import (
	"net/http"
	"testing"

	"github.com/podhmo/go-webtest/try"
)

func TestRequest(t *testing.T) {
	t.Run("response", func(t *testing.T) {
		handler := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"message": "ok", "id": "1"}`))
			},
		)

		try.JSONRequest(
			t,
			handler,
			"GET",
			"/",
			http.StatusOK,
			try.WithAssertJSONResponse(`{"message": "ok", "id": "1"}`),
		)
	})

	t.Run("modifyRequest", func(t *testing.T) {
		handler := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				username, password, ok := r.BasicAuth()
				if !ok {
					t.Fatalf("expecting, basic auth is set. but not set. header %v", r.Header)
				}

				if expected, actual := "*username*", username; expected != actual {
					t.Errorf("invalid username, expected=%q, actual=%q", expected, actual)
				}
				if expected, actual := "*password*", password; expected != actual {
					t.Errorf("invalid password, expected=%q, actual=%q", expected, actual)
				}
			},
		)

		try.JSONRequest(
			t,
			handler,
			"GET",
			"/",
			http.StatusOK,
			try.WithModifyRequest(func(req *http.Request) {
				req.SetBasicAuth("*username*", "*password*")
			}),
		)
	})
	t.Run("assertFunc", func(t *testing.T) {
		handler := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("location", "/newplace")
				w.WriteHeader(http.StatusSeeOther)
			},
		)
		try.JSONRequest(
			t,
			handler,
			"GET",
			"/",
			http.StatusSeeOther,
			try.WithAssertFunc(func(t testing.TB, output *try.Output) {
				location := output.Response.Header.Get("location")
				if expected, actual := "/newplace", location; expected != actual {
					t.Errorf("invalid redirect-location, expected=%q, actual=%q", expected, actual)
				}
			}),
		)
	})
}
