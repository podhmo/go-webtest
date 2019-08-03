package httpbin_test

import (
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/httpbin/httpbintest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/noerror"
)

func TestIt(t *testing.T) {
	ts, teardown := httpbintest.NewTestAPIServer()
	defer teardown()
	client := webtest.NewClientForServer(ts)

	t.Run("200", func(t *testing.T) {
		got, err, teardown := client.Get("/status/200")
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.StatusCode(), err),
			"response: ", got.LazyBodyString(), // add more contextual information?
		)
		defer teardown()

		// todo: assertion response
		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.FromRawWithBytes(got.Data(), got.Body()),
				jsonequal.FromString(`{"message": "OK", "status": 200}`),
			),
		)

		// todo: assertion db check
	})
}

func TestUnit(t *testing.T) {
	handler := httpbintest.NewTestHandler()
	client := webtest.NewClientForHandler(handler)

	t.Run("200", func(t *testing.T) {
		got, err, teardown := client.Get("/status/200")
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.StatusCode(), err),
			"response: ", got.LazyBodyString(), // add more contextual information?
		)
		defer teardown()

		// todo: assertion response
		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.FromRawWithBytes(got.Data(), got.Body()),
				jsonequal.FromString(`{"message": "OK", "status": 200}`),
			),
		)

		// todo: assertion db check
	})
}
