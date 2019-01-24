package try

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

// Input :
type Input struct {
	Method string
	Path   string
	Body   io.Reader

	assertions []func(t testing.TB, res *Output)
	callbacks  []func(req *http.Request)
	bodyString string
}

// Output :
type Output struct {
	Input    *Input
	Body     bytes.Buffer
	Response *http.Response
}

// JSONRequest :
func JSONRequest(t testing.TB, mux http.Handler, method, path string, status int, options ...func(*Input) error) *Output {
	t.Helper()
	return Request(t, mux, method, path, status, append(
		options,
		WithModifyRequest(func(req *http.Request) {
			req.Header.Set("Content-Type", "application/json")
		}),
	)...)
}

// Request :
func Request(t testing.TB, mux http.Handler, method, path string, status int, options ...func(*Input) error) *Output {
	t.Helper()
	input := &Input{
		Method: method,
		Path:   path,
	}

	for _, op := range options {
		if err := op(input); err != nil {
			t.Fatalf("apply option %+v", err)
			return nil
		}
	}
	req := httptest.NewRequest(input.Method, input.Path, input.Body)
	for _, cb := range input.callbacks {
		cb(req)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	res := rec.Result()
	output := Output{
		Input:    input,
		Response: res,
	}

	{
		if _, err := io.Copy(&output.Body, res.Body); err != nil {
			t.Fatalf("parse responose, something wrong: %+v", err)
			return nil
		}
		res.Body = ioutil.NopCloser(&output.Body)
	}

	if expected, got := status, res.StatusCode; got != expected {
		t.Fatalf("status code: expected %d, but %d\n%s", expected, got, output.Body.String())
		return nil
	}

	for _, assert := range input.assertions {
		assert(t, &output)
	}
	return &output
}

// WithJSONBody :
func WithJSONBody(body string) func(input *Input) error {
	return func(input *Input) error {
		input.bodyString = body
		input.Body = bytes.NewBufferString(body)
		return nil
	}
}

// WithModifyRequest :
func WithModifyRequest(callback func(*http.Request)) func(input *Input) error {
	return func(input *Input) error {
		input.callbacks = append(input.callbacks, callback)
		return nil
	}
}

// WithAssertFunc :
func WithAssertFunc(assert func(t testing.TB, output *Output)) func(input *Input) error {
	return func(input *Input) error {
		input.assertions = append(input.assertions, assert)
		return nil
	}
}

// WithAssertJSONResponse :
func WithAssertJSONResponse(body string) func(input *Input) error {
	return func(input *Input) error {
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

		input.assertions = append(input.assertions, func(t testing.TB, output *Output) {
			t.Helper()
			var actual string
			var ob interface{}

			{
				decoder := json.NewDecoder(&output.Body)
				if err := decoder.Decode(&ob); err != nil {
					t.Fatalf("unexpected response:\n%q", output.Body.String())
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
					output.Input.Method,
					output.Input.Path,
					output.Input.bodyString,
					expected,
					actual)
			}
		})
		return nil
	}
}
