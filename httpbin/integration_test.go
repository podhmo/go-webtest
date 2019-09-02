package httpbin_test

import (
	"fmt"
	"net/http"
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
		got, err, teardown := client.GET("/status/200")
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

	t.Run("with hooks", func(t *testing.T) {
		t.Run("200, status check", func(t *testing.T) {
			got, err, teardown := client.GET("/status/200",
				hook.ExpectCode(t, 200),
			)
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
					got, err, teardown := client.GET(c.path,
						hook.GetExpectedDataFromSnapshot(t, &want),
					)
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

	t.Run("auth", func(t *testing.T) {
		cases := []struct {
			user      string
			pass      string
			code      int
			assertion func(*testing.T, webtest.Response)
		}{
			{
				code: 200,
				user: "user", pass: "pass",
				assertion: func(t *testing.T, got webtest.Response) {
					noerror.Should(t,
						jsonequal.ShouldBeSame(
							jsonequal.FromString(`{"authenticated": true, "user": "user"}`),
							jsonequal.From(got.JSONData()),
						),
					)
				},
			},
			{
				code: 401,
				user: "user", pass: "another",
				assertion: func(t *testing.T, got webtest.Response) {
					noerror.Should(t,
						noerror.Equal(401).Actual(got.Code()),
					)
				},
			},
		}

		for _, c := range cases {
			c := c
			t.Run(fmt.Sprintf("%d", c.code), func(t *testing.T) {
				webtest.
					AssertWith(t, c.assertion).
					Try(client.GET("/auth/basic-auth/user/pass",
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
		got, err, teardown := client.Do("GET", "/status/200")
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
					hook.ExpectCode(t, 200),
				}
				if c.query != nil {
					options = append(options, webtest.WithQuery(c.query))
				}

				got, _, teardown := client.Do("GET", c.path, options...)
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
