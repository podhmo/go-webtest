package webtest

import (
	"testing"
)

// AssertWith :
func AssertWith(t *testing.T, assertion func(t *testing.T, got Response)) *TryWithAssertion {
	return &TryWithAssertion{
		t:         t,
		assertion: assertion,
	}
}

// TryWithAssertion :
type TryWithAssertion struct {
	t         *testing.T
	assertion func(*testing.T, Response)
}

// Try :
func (a *TryWithAssertion) Try(got Response, err error, teardown func()) {
	if err != nil {
		a.t.Fatalf("try: %+v", err)
	}
	defer teardown()
	a.assertion(a.t, got)
}
