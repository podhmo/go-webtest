package httpbin_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/podhmo/go-webtest/httpbin/httpbintest"
	"github.com/podhmo/go-webtest/jsonequal"
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

		// todo: assertion response
		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.FromReadCloser(res.Body),
				jsonequal.FromString(`{"message": "OK", "status": 200}`),
			),
		)

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

		// todo: assertion response
		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.FromReadCloser(res.Body),
				jsonequal.FromString(`{"message": "OK", "status": 200}`),
			),
		)

		// todo: assertion db check
	})
}
