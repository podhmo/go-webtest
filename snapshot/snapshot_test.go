package snapshot

import (
	"encoding/json"
	"os"
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

func TestWithMetadata(t *testing.T) {
	got := map[string]string{"message": "hello"}
	metadata := map[string]interface{}{
		"path":   "/greeting",
		"method": "GET",
		"status": 200,
	}

	// checking by json
	if want := Take(t, got, WithMetadata(metadata)); !reflect.DeepEqual(normalize(got), normalize(want)) {
		t.Errorf("want %v, but got %v", want, got)
	}

	// check stored file
	{
		f, err := os.Open(NewTestdataRecorder(nil).Path(t))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		decoder := json.NewDecoder(f)
		storedData := loadData{}
		if err := decoder.Decode(&storedData); err != nil {
			t.Fatal(err)
		}
		if storedData.Metadata == nil {
			t.Fatal("metadata must be existed")
		}

		if !reflect.DeepEqual(normalize(metadata), normalize(storedData.Metadata)) {
			t.Errorf("metadata, want %v, but %v", metadata, storedData.Metadata)
		}
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
