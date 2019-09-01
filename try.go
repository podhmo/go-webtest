package webtest

import (
	"testing"
)

// AssertWith :
func AssertWith(t testing.TB, assertion func(got Response)) *TryWithAssertion {
	return &TryWithAssertion{
		t:         t,
		assertion: assertion,
	}
}

// TryWithAssertion :
type TryWithAssertion struct {
	t         testing.TB
	assertion func(Response)
}

// Try :
func (a *TryWithAssertion) Try(got Response, err error, teardown func()) {
	if err != nil {
		a.t.Fatalf("try: %+v", err)
	}
	defer teardown()
	a.assertion(got)
}
