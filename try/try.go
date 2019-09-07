package try

import (
	"testing"

	webtest "github.com/podhmo/go-webtest"
	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/tripperware"
	"github.com/podhmo/noerror"
)

// It :
type It struct {
	_Required_Sentinel struct{}

	Code           int
	Want           *interface{}
	ModifyResponse func(res webtest.Response) (got interface{})
}

// With :
func (it It) With(
	t *testing.T,
	c *webtest.Client,
	method string,
	path string,
	options ...webtest.Option,
) {
	t.Helper()
	needSnapshot := it.Want != nil

	if it.Code != 0 {
		options = append(options, webtest.WithTripperware(tripperware.ExpectCode(t, it.Code)))
	}
	if needSnapshot {
		options = append(options, webtest.WithTripperware(tripperware.GetExpectedDataFromSnapshot(t, it.Want)))
	}

	got, err := c.Do(method, path, options...)
	noerror.Must(t, err)
	defer func() { noerror.Must(t, got.Close()) }()

	if needSnapshot {
		noerror.Must(t, noerror.NotEqual(nil).Actual(it.Want))

		var actual interface{}
		if it.ModifyResponse != nil {
			actual = it.ModifyResponse(got)
		} else {
			actual = got.JSONData()
		}

		mismatch := jsonequal.ShouldBeSame(
			jsonequal.From(actual),
			jsonequal.From(it.Want),
			jsonequal.WithLeftName("left (got) "),
			jsonequal.WithRightName("right (want) "),
		)
		if mismatch != nil {
			t.Errorf("%s", mismatch)
		}
	}
}
