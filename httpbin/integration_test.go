package httpbin_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/httpbin/httpbintest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/tripperware"
	"github.com/podhmo/noerror"
)

func TestIt(t *testing.T) {
	ts, teardown := httpbintest.NewTestAPIServer()
	defer teardown()
	client := webtest.NewClientFromTestServer(ts)

	t.Run("200", func(t *testing.T) {
		got, err := client.Get("/status/200")
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.Code(), err),
			"response: ", got.LazyText(), // add more contextual information?
		)
		defer func() { noerror.Should(t, got.Close()) }()

		// todo: assertion response
		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.FromRawWithBytes(got.JSONData(), got.Body()),
				jsonequal.FromString(`{"message": "OK", "status": 200}`),
			),
		)

		// todo: assertion db check
	})

	t.Run("with hooks", func(t *testing.T) {
		t.Run("200, status check", func(t *testing.T) {
			got, err := client.Get("/status/200",
				webtest.WithTripperware(
					tripperware.ExpectCode(t, 200),
				),
			)
			noerror.Must(t, err)
			defer func() { noerror.Should(t, got.Close()) }()

			noerror.Should(t,
				jsonequal.ShouldBeSame(
					jsonequal.FromRawWithBytes(got.JSONData(), got.Body()),
					jsonequal.FromString(`{"message": "OK", "status": 200}`),
				),
			)
		})

		t.Run("snapshot", func(t *testing.T) {
			var want interface{}

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
					got, err := client.Get(c.path,
						webtest.WithTripperware(
							tripperware.GetExpectedDataFromSnapshot(t, &want),
						),
					)
					noerror.Must(t, err)
					defer func() { noerror.Should(t, got.Close()) }()

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

	t.Run("auth", func(t *testing.T) {
		cases := []struct {
			user      string
			pass      string
			code      int
			assertion webtest.Assertion
		}{
			{
				code: 200,
				user: "user", pass: "pass",
				assertion: func(t *testing.T, got webtest.Response) {
					noerror.Should(t,
						jsonequal.ShouldBeSame(
							jsonequal.From(got.JSONData()),
							jsonequal.FromString(`{"authenticated": true, "user": "user"}`),
						),
					)
				},
			},
			{
				code: 401,
				user: "user", pass: "another",
			},
		}

		for _, c := range cases {
			c := c
			t.Run(fmt.Sprintf("%d", c.code), func(t *testing.T) {
				webtest.
					Try(t, c.assertion).
					With(client.Get("/auth/basic-auth/user/pass",
						webtest.WithTripperware(
							tripperware.ExpectCode(t, c.code),
						),
						webtest.WithModifyRequest(func(req *http.Request) {
							req.SetBasicAuth(c.user, c.pass)
						}),
					))
			})
		}
	})
}

func TestUnit(t *testing.T) {
	handler := httpbintest.NewTestHandler()
	client := webtest.NewClientFromHandler(handler)

	t.Run("200", func(t *testing.T) {
		got, err := client.Do("GET", "/status/200")
		noerror.Must(t, err)
		noerror.Must(t,
			noerror.Equal(200).ActualWithError(got.Code(), err),
			"response: ", got.LazyText(), // add more contextual information?
		)
		defer func() { noerror.Should(t, got.Close()) }()

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
				options := []webtest.Option{
					webtest.WithTripperware(
						tripperware.ExpectCode(t, 200),
					),
				}
				if c.query != nil {
					options = append(options, webtest.WithQuery(c.query))
				}

				got, err := client.Do("GET", c.path, options...)
				noerror.Must(t, err)
				defer func() { noerror.Should(t, got.Close()) }()

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
