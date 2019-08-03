package httpbin_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/podhmo/go-webtest/httpbin/httpbintest"
	"github.com/podhmo/noerror"
)

func TestIt(t *testing.T) {
	ts, teardown := httpbintest.NewTestAPIServer()
	defer teardown()

	t.Run("200", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/status/200", ts.URL))

		noerror.Must(t,
			noerror.Equal(200).ActualWithError(res.StatusCode, err),
			"response: ", "<todo>", // add more contextual information?
		)

		data := map[string]interface{}{}
		decoder := json.NewDecoder(res.Body)
		noerror.Must(t,
			decoder.Decode(&data),
			"response: ", "<todo>",
		)
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

		noerror.Must(t,
			noerror.Equal(200).Actual(res.StatusCode),
			"response: ", "<todo>",
		)

		data := map[string]interface{}{}
		decoder := json.NewDecoder(res.Body)
		noerror.Must(t,
			decoder.Decode(&data),
			"response: ", "<todo>",
		)
		defer res.Body.Close()

		// todo: assertion response
		fmt.Printf("body: %#+v", data)

		// todo: assertion db check
	})
}
