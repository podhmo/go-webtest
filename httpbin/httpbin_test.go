package httpbin_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/podhmo/go-webtest/httpbin/httpbintest"
)

func TestIt(t *testing.T) {
	ts, teardown := httpbintest.NewTestAPIServer()
	defer teardown()

	t.Run("200", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/status/200", ts.URL))
		if err != nil {
			t.Fatalf("%+v", err) // add more contextual information?
		}
		defer teardown()

		if res.StatusCode != 200 {
			t.Fatalf("status expect 200, but %d\n response: %s", res.StatusCode, "<todo>")
		}

		data := map[string]interface{}{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&data); err != nil {
			t.Fatalf("parse error %+v\n response:%s", err, "<todo>")
		}
		defer res.Body.Close()

		// todo: assertion response
		fmt.Printf("body: %#+v", data)

		// todo: assertion db check
	})
}

func TestUnit(t *testing.T) {
	handler := httpbintest.NewTestHandler()

	t.Run("200", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/status/200", nil)
		handler(w, req)
		res := w.Result()
		if res.StatusCode != 200 {
			t.Fatalf("status expect 200, but %d\n response: %s", res.StatusCode, "<todo>")
		}

		data := map[string]interface{}{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&data); err != nil {
			t.Fatalf("parse error %+v\n response:%s", err, "<todo>")
		}
		defer res.Body.Close()

		// todo: assertion response
		fmt.Printf("body: %#+v", data)

		// todo: assertion db check
	})
}
