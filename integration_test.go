package webtest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/tripperware"
	"github.com/podhmo/noerror"
)

type Input struct {
	Values []int
}

func Add(w http.ResponseWriter, req *http.Request) {
	var data Input
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&data); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"error": %q}`, err.Error())
		return
	}

	var n int
	for _, v := range data.Values {
		n += v
	}
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(map[string]int{"result": n}); err != nil {
		panic(err)
	}
}

func TestHandler(t *testing.T) {
	c := webtest.NewClientFromHandler(http.HandlerFunc(Add))
	var want interface{}
	got, err := c.Post("/",
		webtest.WithJSON(bytes.NewBufferString(`{"values": [1,2,3]}`)),
		webtest.WithTripperware(
			tripperware.ExpectCode(t, 200),
			tripperware.GetExpectedDataFromSnapshot(t, &want),
		),
	)

	noerror.Must(t, err)
	noerror.Should(t,
		jsonequal.ShouldBeSame(
			jsonequal.From(got.JSONData()),
			jsonequal.From(want),
		),
	)
}
