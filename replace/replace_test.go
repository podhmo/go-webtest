package replace

import (
	"testing"

	"github.com/podhmo/go-webtest/jsonequal"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestByMap(t *testing.T) {
	dummy := &person{
		Name: "foo",
		Age:  10,
	}

	// store data's age is 20, but actual's one is 10
	want := []byte(`{"name": "foo", "age": 20}`)

	data := jsonequal.MustNormalize(dummy)
	refMap := map[string]interface{}{
		"/age": 20, // replace 10 -> 20
	}
	if _, err := ByMap(data, refMap); err != nil {
		t.Error(err)
	}
	if err := jsonequal.ShouldBeSame(jsonequal.From(data), jsonequal.FromBytes(want)); err != nil {
		t.Error(err)
	}
}

func TestByMap2(t *testing.T) {
	data, _, _ := jsonequal.FromString(`{"a": {"b": {"c": {"target": 100}}}}`)()
	want := []byte(`{"a": {"b": {"c": {"target": 200}}}}`)

	refMap := map[string]interface{}{
		"#/a/b/c/target": 200,
	}
	if _, err := ByMap(data, refMap); err != nil {
		t.Error(err)
	}
	if err := jsonequal.ShouldBeSame(jsonequal.From(data), jsonequal.FromBytes(want)); err != nil {
		t.Error(err)
	}
}

func TestByPalette(t *testing.T) {
	dummy := &person{Name: "foo", Age: 10}
	want := []byte(`{"name": "foo", "age": 20}`)

	data := jsonequal.MustNormalize(dummy)
	refs := []string{"/age"}
	pallete := jsonequal.MustNormalize(person{Age: 20})

	if _, err := ByPalette(data, refs, pallete); err != nil {
		t.Error(err)
	}
	if err := jsonequal.ShouldBeSame(jsonequal.From(data), jsonequal.FromBytes(want)); err != nil {
		t.Error(err)
	}
}
