package try_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/podhmo/go-webtest/try"
)

// helper for testing

type call struct {
	Method string
	Format string
	Args   []interface{}
}

func (c call) String() string {
	return c.Method + ":" + fmt.Sprintf(c.Format, c.Args...)
}

type fakeT struct {
	*testing.T
	Called []call
}

// Errorf :
func (ft *fakeT) Errorf(fmt string, args ...interface{}) {
	ft.Called = append(ft.Called, call{
		Method: "Errorf",
		Format: fmt,
		Args:   args,
	})
}

// Fatalf :
func (ft *fakeT) Fatalf(fmt string, args ...interface{}) {
	ft.Called = append(ft.Called, call{
		Method: "Fatalf",
		Format: fmt,
		Args:   args,
	})
}

func TestRequestMismatch(t *testing.T) {
	cases := []struct {
		Name         string
		MustIncluded string
		Handler      http.HandlerFunc
		Options      []func(*try.Input) error
	}{
		{
			Name:         "status",
			MustIncluded: `{"message": "not found"}`,
			Handler: http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					io.WriteString(w, `{"message": "not found"}`)
				},
			),
		},
		{
			Name:         "response",
			MustIncluded: `"message": "ok"`,
			Handler: http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					io.WriteString(w, `{"message": "ok"}`)
				},
			),
			Options: []func(*try.Input) error{
				try.WithAssertJSONResponse(`{"status": "ok"}`),
			},
		},
		{
			Name:         "response2",
			MustIncluded: `"status": "ok"`,
			Handler: http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					io.WriteString(w, `{"message": "ok"}`)
				},
			),
			Options: []func(*try.Input) error{
				try.WithAssertJSONResponse(`{"status": "ok"}`),
			},
		},
		{
			Name:         "path",
			MustIncluded: `GET /`,
			Handler: http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					io.WriteString(w, `{"message": "ok"}`)
				},
			),
			Options: []func(*try.Input) error{
				try.WithAssertJSONResponse(`{"status": "ok"}`),
			},
		},
		{
			Name:         "diff",
			MustIncluded: `## diff (- missing, + excess)`,
			Handler: http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					io.WriteString(w, `{"message": "ok"}`)
				},
			),
			Options: []func(*try.Input) error{
				try.WithAssertJSONResponse(`{"status": "ok"}`),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			ft := &fakeT{T: t}

			// must error
			try.Request(ft, c.Handler, "GET", "/", http.StatusOK, c.Options...)

			if expected := 1; len(ft.Called) != expected {
				t.Fatalf("unexpected calling, expected call count is %d, but actual is %d", expected, len(ft.Called))
			}

			{
				expected := c.MustIncluded
				actual := ft.Called[0].String()
				if !strings.Contains(actual, expected) {
					t.Errorf("expecting message is not included.\nexpect:\n%s\nactual:\n%s", expected, actual)
				}
			}
		})
	}
}
