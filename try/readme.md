## go-webtest/try

install

```console
$ go get -v github.com/podhmo/go-webtest/try
```

handler

```go
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
```

test output

```
--- FAIL: Test (0.00s)
    --- FAIL: Test/status_mismatch (0.00s)
        example_test.go:30: status code: expected 200, but 400
            {"error": "invalid value"}
    --- FAIL: Test/response_mismatch (0.00s)
        example_test.go:43: mismatch response:
            ## diff (- missing, + excess)
            
              {
            -   "message": "hello WORLD",
            -   "name": "WORLD"
            +   "message": "hello world"
              }
            
            ## request
            
            GET /
            
            {"name": "WORLD"}
            
            ## expected response
            
            {
              "message": "hello world"
            }
            
            ## actual response
            {
              "message": "hello WORLD",
              "name": "WORLD"
            }
FAIL
exit status 1
FAIL	github.com/podhmo/go-webtest/_examples/try	0.005s
```

test code

```go
package examples

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/podhmo/go-webtest/try"
)

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
```


