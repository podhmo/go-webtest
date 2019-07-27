package snapshot

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestSimple(t *testing.T) {
	got := GetData() // see: helper_test.go

	// checking by json
	if want := Take(t, got); !reflect.DeepEqual(normalize(got), normalize(want)) {
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
