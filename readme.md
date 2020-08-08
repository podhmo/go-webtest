[![CircleCI](https://circleci.com/gh/podhmo/go-webtest.svg?style=svg)](https://circleci.com/gh/podhmo/go-webtest)

# go-webtest

features

- handling response by [custom interface](https://godoc.org/github.com/podhmo/go-webtest/testclient#Response) 
- debug tracing when `DEBUG=1`
- snapshot testing (if update snapshot, `SNAPSHOT=1` or `SNAPSHOT=<golden file path>`)
- json diff with https://github.com/nsf/jsondiff

## examples

full test code is [here](./integration_test.go), the test target handler is defined [here](https://github.com/podhmo/go-webtest/blob/10353a9f1e700503028b420bf3068781030e5dac/integration_test.go#L25-L43)

### with webtest

```go
import (
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/tripperware"
	"github.com/podhmo/noerror"
)

func TestWithWebtest(t *testing.T) {
	c := webtest.NewClientFromHandler(http.HandlerFunc(handleAdd))
	var want interface{}
	got, err := c.Post("/",
		webtest.WithJSON(bytes.NewBufferString(`{"values": [1,2,3]}`)),
		webtest.WithTripperware(
			tripperware.ExpectCode(t, 200),
			tripperware.GetExpectedDataFromSnapshot(t, &want),
		),
	)

	noerror.Must(t, err)
	defer func() { noerror.Must(t, got.Close()) }()
	noerror.Should(t,
		jsonequal.ShouldBeSame(
			jsonequal.From(got.JSONData()),
			jsonequal.From(want),
		),
	)
}
```

### with try package (shortcut)

```go
import (
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/try"
)


func TestWithTry(t *testing.T) {
	c := webtest.NewClientFromHandler(http.HandlerFunc(handleAdd))

	var want interface{}
	try.It{
		Code: 200,
		Want: &want,
		ModifyResponse: func(res webtest.Response) (got interface{}) {
			return res.JSONData()
		},
	}.With(t, c,
		"POST", "/",
		webtest.WithJSON(bytes.NewBufferString(`{"values": [1,2,3]}`)),
	)
}
```

If modify response is not needed, it is also ok, when the response does not include *semi-random value* (for example the value of now time).

```go
c := webtest.NewClientFromHandler(http.HandlerFunc(handleAdd))

var want interface{}
try.It{
	Code: 200,
	Want: &want,
}.With(t, c,
	"POST", "/",
	webtest.WithJSON(bytes.NewBufferString(`{"values": [1,2,3]}`)),
)
```

### without-webtest (but snapshot testing)

```go
import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/podhmo/go-webtest/snapshot"
)

func TestWithoutWebtest(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/", bytes.NewBufferString(`{"values": [1,2,3]}`))
	req.Header.Set("Content-Type", "application/json")

	handleAdd(rec, req)
	res := rec.Result()

	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		t.Fatalf("status code, want 200, but got %d\n response:%s", res.StatusCode, string(b))
	}

	var got interface{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&got); err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	want := snapshot.Take(t, &got)
	if !reflect.DeepEqual(want, got) {
		t.Errorf(`want %s, but got %s`, want, got)
	}
}
```

## the location of snapshot data

The snapshot data is saved in `testdata/<test function name>.golden` (e.g. testdata/TestHandler/try.golden) .

```json
{
  "modifiedAt": "2019-09-07T21:40:30.70331035+09:00",
  "data": {
    "request": {
      "method": "POST",
      "path": "/"
    },
    "response": {
      "data": {
        "result": 6
      },
      "statusCode": 200
    }
  }
}
```

## example output if tests are failed

Output examples.

### ✅ debug trace (not failed)

```console
$ DEBUG=1 go test -v
2019/09/08 08:42:56 builtin debug trace is activated
=== RUN   TestHandler
=== RUN   TestHandler/plain
=== RUN   TestHandler/webtest
	Request : ------------------------------
	POST / HTTP/1.1
	Host: example.com
	Content-Type: application/json
	
	{"values": [1,2,3]}
	----------------------------------------
	Response: ------------------------------
	HTTP/1.1 200 OK
	Connection: close
	Content-Type: application/json
	
	{"result":6}
	----------------------------------------
=== RUN   TestHandler/try
	Request : ------------------------------
	POST / HTTP/1.1
	Host: example.com
	Content-Type: application/json
	
	{"values": [1,2,3]}
	----------------------------------------
	Response: ------------------------------
	HTTP/1.1 200 OK
	Connection: close
	Content-Type: application/json
	
	{"result":6}
	----------------------------------------
--- PASS: TestHandler (0.00s)
    --- PASS: TestHandler/plain (0.00s)
        snapshot.go:56: load snapshot data: "testdata/TestHandler/plain.golden"
    --- PASS: TestHandler/webtest (0.00s)
        snapshot.go:56: load snapshot data: "testdata/TestHandler/webtest.golden"
    --- PASS: TestHandler/try (0.00s)
        snapshot.go:56: load snapshot data: "testdata/TestHandler/try.golden"
PASS
ok  	github.com/podhmo/go-webtest	0.005s
```

### ❌ unexpected status

```console
$ go test
--- FAIL: TestHandler (0.00s)
    --- FAIL: TestHandler/try (0.00s)
        snapshot.go:56: load snapshot data: "testdata/TestHandler/try.golden"
        integration_test.go:103: unexpected error, status code, got is "200 OK", but want is "201 Created"
            Response: ------------------------------
            HTTP/1.1 200 OK
            Connection: close
            Content-Type: application/json
            
            {"result":6}
            ----------------------------------------
            
FAIL
```

### ❌ unexpected response

```console
$ go test
--- FAIL: TestHandler (0.00s)
    --- FAIL: TestHandler/try (0.00s)
        snapshot.go:56: load snapshot data: "testdata/TestHandler/try.golden"
        integration_test.go:103: on equal check: jsondiff, got and want is not same.
            status=NoMatch
            {
                "result": 6 => 10
            }
            
            left (want) :
            	{"result":6}
            right (got) :
            	{"result":10}
FAIL
```

## sub packages

### try

todo

### snapshot

todo

#### snapshot/replace

todo

### jsonequal

todo

### testclient

todo
