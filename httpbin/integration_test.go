package httpbin_test

import (
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/httpbin/httpbintest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/middlewares"
	"github.com/podhmo/noerror"
)

func TestIt(t *testing.T) {
	ts, teardown := httpbintest.NewTestAPIServer()
	defer teardown()
	client := webtest.NewClientFromTestServer(ts)

	t.Run("200", func(t *testing.T) {
		got, err, teardown := client.Do(t, "/status/200")
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.StatusCode(), err),
			"response: ", got.LazyBodyString(), // add more contextual information?
		)
		defer teardown()

		// todo: assertion response
		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.FromRawWithBytes(got.JSONData(), got.Body()),
				jsonequal.FromString(`{"message": "OK", "status": 200}`),
			),
		)

		// todo: assertion db check
	})

	t.Run("with middlewares", func(t *testing.T) {
		t.Run("200, status check", func(t *testing.T) {
			client := client.Bind(
				middlewares.ExpectStatusCode(200),
			)
			got, err, teardown := client.Do(t, "/status/200")
			noerror.Must(t, err)
			defer teardown()

			noerror.Should(t,
				jsonequal.ShouldBeSame(
					jsonequal.FromRawWithBytes(got.JSONData(), got.Body()),
					jsonequal.FromString(`{"message": "OK", "status": 200}`),
				),
			)
		})

		t.Run("snapshot", func(t *testing.T) {
			var want interface{}

			client := client.Bind(
				middlewares.SnapshotTesting(&want),
			)

			cases := []struct {
				path string
				msg  string
			}{
				{msg: "200", path: "/status/200"},
				{msg: "201", path: "/status/201"},
				{msg: "404", path: "/status/404"},
			}

			for _, c := range cases {
				c := c
				t.Run(c.msg, func(t *testing.T) {
					got, err, teardown := client.Do(t, c.path)
					noerror.Must(t, err)
					defer teardown()

					noerror.Should(t,
						jsonequal.ShouldBeSame(
							jsonequal.FromRaw(want),
							jsonequal.FromRawWithBytes(got.JSONData(), got.Body()),
						),
					)
				})
			}
		})
	})
}

func TestUnit(t *testing.T) {
	handler := httpbintest.NewTestHandler()
	client := webtest.NewClientFromHandler(handler)

	t.Run("200", func(t *testing.T) {
		got, err, teardown := client.Do(t, "/status/200")
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.StatusCode(), err),
			"response: ", got.LazyBodyString(), // add more contextual information?
		)
		defer teardown()

		// todo: assertion response
		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.FromRawWithBytes(got.JSONData(), got.Body()),
				jsonequal.FromString(`{"message": "OK", "status": 200}`),
			),
		)

		// todo: assertion db check
	})
}
