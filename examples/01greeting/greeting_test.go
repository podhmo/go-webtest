package greeting

import (
	"net/http"
	"testing"
	"time"

	"github.com/podhmo/go-rfc3339"
	"github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/tripperware"
	"github.com/podhmo/noerror"
)

func TestIt(t *testing.T) {
	app := &App{
		Clock: ClockFunc(func() time.Time {
			return rfc3339.MustParse("2000-01-01T00:00:00Z")
		}),
	}

	var want interface{}
	client := webtest.NewClientFromHandler(http.HandlerFunc(app.Greeting))
	got, err := client.Get("/",
		webtest.WithTripperware(
			tripperware.ExpectCode(t, 200),
			tripperware.GetExpectedDataFromSnapshot(t, &want),
		),
	)
	noerror.Must(t, err)
	noerror.Should(t,
		jsonequal.ShouldBeSame(
			jsonequal.FromRawWithBytes(got.JSONData(), got.Body()),
			jsonequal.FromRaw(want),
		),
	)
}
