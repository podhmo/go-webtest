package webtest

import (
	"testing"
)

// Assertion :
type Assertion = func(t *testing.T, got Response)

// AssertWith :
func AssertWith(t *testing.T, assertions ...Assertion) *TryWithAssertion {
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

// Try :
func (a *TryWithAssertion) Try(got Response, err error, teardown func()) {
	if err != nil {
		a.t.Fatalf("try: %+v", err)
	}
	defer teardown()
	for _, assert := range a.assertions {
		assert(a.t, got)
	}
}
