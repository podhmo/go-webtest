package httpbin_test

import (
	"fmt"
	"net/url"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/hook"
	"github.com/podhmo/go-webtest/httpbin/httpbintest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/noerror"
)

func TestIt(t *testing.T) {
	ts, teardown := httpbintest.NewTestAPIServer()
	defer teardown()
	client := webtest.NewClientFromTestServer(ts)

	t.Run("200", func(t *testing.T) {
		got, err, teardown := client.GET(t, "/status/200")
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.Code(), err),
			"response: ", got.LazyText(), // add more contextual information?
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
				hook.ExpectCode(200),
			)
			got, err, teardown := client.GET(t, "/status/200")
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
				hook.GetExpectedDataFromSnapshot(&want),
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
					got, err, teardown := client.GET(t, c.path)
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
		got, err, teardown := client.Do(t, "GET", "/status/200")
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.Code(), err),
			"response: ", got.LazyText(), // add more contextual information?
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

	t.Run("get", func(t *testing.T) {
		client = client.Bind(
			hook.ExpectCode(200),
		)

		cases := []struct {
			path     string
			query    url.Values
			expected interface{}
		}{
			{
				path:     "/get",
				expected: map[string][]string{},
			},
			{
				path:     "/get?xxx=111",
				expected: map[string][]string{"xxx": []string{"111"}},
			},
			{
				path:  "/get?xxx=111",
				query: webtest.MustParseQuery("yyy=222"),
				expected: map[string][]string{
					"xxx": []string{"111"},
					"yyy": []string{"222"},
				},
			},
			{
				path:  "/get?xxx=111",
				query: webtest.MustParseQuery("yyy=222&xxx=333"),
				expected: map[string][]string{
					"xxx": []string{"333", "111"},
					"yyy": []string{"222"},
				},
			},
		}
		for i, c := range cases {
			c := c
			t.Run(fmt.Sprintf("case%d", i), func(t *testing.T) {
				var options []func(*webtest.Config)
				if c.query != nil {
					options = append(options, webtest.WithQuery(c.query))
				}
				got, _, teardown := client.Do(t, "GET", c.path, options...)
				defer teardown()

				var data map[string]interface{}
				noerror.Must(t, got.ParseJSONData(&data))

				noerror.Should(t,
					jsonequal.ShouldBeSame(
						jsonequal.From(data["args"]),
						jsonequal.From(c.expected),
					),
				)
			})
		}
	})
}
