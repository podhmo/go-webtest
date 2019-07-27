package snapshot_test

import (
	"testing"

	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/snapshot"
)

func TestIt(t *testing.T) {
	got := snapshot.GetData() // see: helper_test.go
	want := snapshot.Take(t, got)
	if err := jsonequal.ShouldBeSame(jsonequal.From(got), jsonequal.From(want)); err != nil {
		t.Errorf("%+v", err)
	}
}
