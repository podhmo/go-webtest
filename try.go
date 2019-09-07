package webtest

import (
	"testing"

	"github.com/podhmo/noerror"
)

// Assertion :
type Assertion = func(t *testing.T, got Response)

// Try :
func Try(t *testing.T, assertions ...Assertion) *TryWithAssertion {
	var args []Assertion
	for _, arg := range assertions {
		if arg == nil {
			continue
		}
		args = append(args, arg)
	}
	return &TryWithAssertion{
		t:          t,
		assertions: args,
	}
}

// TryWithAssertion :
type TryWithAssertion struct {
	t          *testing.T
	assertions []Assertion
}

// With :
func (a *TryWithAssertion) With(got Response, err error) {
	a.t.Helper()
	noerror.Must(a.t, err)
	defer func() { noerror.Should(a.t, got.Close()) }()

	for _, assert := range a.assertions {
		assert(a.t, got)
	}
}
