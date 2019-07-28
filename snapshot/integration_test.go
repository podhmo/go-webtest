package snapshot_test

import (
	"testing"
	"time"

	rfc3339 "github.com/podhmo/go-rfc3339"
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

func TestItWithReplaceMap(t *testing.T) {
	got := map[string]string{
		"name": "foo",
		"now":  rfc3339.Format(time.Now()),
	}
	repMap := map[string]interface{}{
		"#/now": got["now"],
	}
	want := snapshot.Take(t, got, snapshot.WithReplaceMap(repMap))

	if err := jsonequal.ShouldBeSame(jsonequal.From(got), jsonequal.From(want)); err != nil {
		t.Errorf("%+v", err)
	}
}
