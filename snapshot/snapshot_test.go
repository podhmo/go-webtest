package snapshot_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/podhmo/go-webtest/snapshot"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func getData() person {
	return person{
		Name: "foo",
		Age:  20,
	}
}

func TestSimple(t *testing.T) {
	// checking by json
	got := getData()
	if want := snapshot.Take(t, got); !reflect.DeepEqual(normalize(got), normalize(want)) {
		t.Errorf("want %v, but got %v", want, got)
	}
}

func normalize(src interface{}) interface{} {
	b, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	var dst interface{}
	if err := json.Unmarshal(b, &dst); err != nil {
		panic(err)
	}
	return dst
}
