package webtest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cside/jsondiff"
	"github.com/pkg/errors"
)

// TryRequestRequest :
type TryRequestRequest struct {
	Method     string
	Path       string
	Body       io.Reader
	Assertions []func(t testing.TB, res *TryRequestResponse)
	Response   TryRequestResponse

	bodyString string
}

// TryRequestResponse :
type TryRequestResponse struct {
	Request *TryRequestRequest
	Body    bytes.Buffer
	*http.Response
}

// TryRequest :
func TryRequest(t testing.TB, mux http.Handler, method, path string, status int, options ...func(*TryRequestRequest) error) *TryRequestResponse {
	t.Helper()
	treq := &TryRequestRequest{
		Method: method,
		Path:   path,
	}

	for _, op := range options {
		if err := op(treq); err != nil {
			t.Fatalf("apply option %+v", err)
			return nil
		}
	}
	req := httptest.NewRequest(treq.Method, treq.Path, treq.Body)

	// todo: to option
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	res := rec.Result()
	tresponse := TryRequestResponse{
		Request:  treq,
		Response: res,
	}

	{
		if _, err := io.Copy(&tresponse.Body, res.Body); err != nil {
			t.Fatalf("parse responose, something wrong: %+v", err)
			return nil
		}
		res.Body = ioutil.NopCloser(&tresponse.Body)
	}

	if expected, got := status, res.StatusCode; got != expected {
		t.Fatalf("status code: expected %d, but %d\n%s", expected, got, tresponse.Body.String())
		return nil
	}

	for _, assert := range treq.Assertions {
		assert(t, &tresponse)
	}
	return &tresponse
}

// WithJSONBody :
func WithJSONBody(body string) func(treq *TryRequestRequest) error {
	return func(treq *TryRequestRequest) error {
		treq.bodyString = body
		treq.Body = bytes.NewBufferString(body)
		return nil
	}
}

// WithAssert :
func WithAssert(assert func(t testing.TB, res *TryRequestResponse)) func(treq *TryRequestRequest) error {
	return func(treq *TryRequestRequest) error {
		treq.Assertions = append(treq.Assertions, assert)
		return nil
	}
}

// WithAssertJSONResponse :
func WithAssertJSONResponse(body string) func(treq *TryRequestRequest) error {
	return func(treq *TryRequestRequest) error {
		var expected string
		{
			var ob interface{}
			if err := json.Unmarshal([]byte(body), &ob); err != nil {
				return errors.Wrap(err, "prepare unmarsal")
			}
			b, err := json.MarshalIndent(&ob, "", "  ")
			if err != nil {
				return errors.Wrap(err, "prepare marsal")
			}
			expected = string(b)
		}

		treq.Assertions = append(treq.Assertions, func(t testing.TB, res *TryRequestResponse) {
			var actual string
			var ob interface{}

			{
				decoder := json.NewDecoder(&res.Body)
				if err := decoder.Decode(&ob); err != nil {
					t.Fatalf("unexpected response:\n%q", res.Body.String())
				}
				b, err := json.MarshalIndent(&ob, "", "  ")
				if err != nil {
					panic(err) // something wrong
				}
				actual = string(b)
			}
			diff := jsondiff.LineDiff(actual, expected)
			if diff != "" {
				t.Fatalf(`mismatch response:
## diff (- missing, + excess)

%s

## request

%s %s

%s

## expected response

%s

## actual response
%s`,
					diff,
					res.Request.Method,
					res.Request.Path,
					res.Request.bodyString,
					expected,
					actual)
			}
		})
		return nil
	}
}
