package webtest

import (
	"testing"
)

// Assertion :
type Assertion = func(t *testing.T, got Response)

// Try :
func Try(t *testing.T, assertions ...Assertion) *TryWithAssertion {
	return &TryWithAssertion{
		t:          t,
		assertions: assertions,
	}
}

// TryWithAssertion :
type TryWithAssertion struct {
	t          *testing.T
	assertions []Assertion
}

// With :
func (a *TryWithAssertion) With(got Response, err error, teardown func()) {
	a.t.Helper()
	if err != nil {
		a.t.Fatalf("try: %+v", err)
	}
	defer teardown()
	for _, assert := range a.assertions {
		assert(a.t, got)
	}
}
