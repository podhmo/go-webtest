package replace_test

import (
	"testing"
	"time"

	rfc3339 "github.com/podhmo/go-rfc3339"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/snapshot/replace"
)

type Data struct {
	Name string    `json:"name"`
	Now  time.Time `json:"now"`
}

func TestByMap(t *testing.T) {
	now := rfc3339.MustParse("2000-01-01T00:00:00Z")

	dummy := &Data{Name: "foo", Now: time.Now()}
	want := []byte(`{"name": "foo", "now": "2000-01-01T00:00:00Z"}`)

	data := jsonequal.MustNormalize(dummy)
	refMap := map[string]interface{}{
		"/now": now,
	}

	if _, err := replace.ByMap(data, refMap); err != nil {
		t.Error(err)
	}
	if err := jsonequal.ShouldBeSame(jsonequal.From(data), jsonequal.FromBytes(want)); err != nil {
		t.Error(err)
	}
}
