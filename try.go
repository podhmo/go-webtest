package webtest

import (
	"testing"

	"github.com/podhmo/go-webtest/jsonequal"
	"github.com/podhmo/go-webtest/tripperware"
	"github.com/podhmo/noerror"
)

// Try :
type Try struct {
	_Required_Sentinel struct{}

	Code int

	Data           *interface{}
	ModifyResponse func(res Response) (actual interface{})
}

// Do :
func (expected Try) Do(
	t *testing.T,
	c *Client,
	method string,
	path string,
	options ...Option,
) {
	t.Helper()
	needSnapshot := expected.Data != nil

	if expected.Code != 0 {
		options = append(options, WithTripperware(tripperware.ExpectCode(t, expected.Code)))
	}
	if needSnapshot {
		options = append(options, WithTripperware(tripperware.GetExpectedDataFromSnapshot(t, expected.Data)))
	}

	got, err := c.Do(method, path, options...)
	noerror.Must(t, err)
	defer func() { noerror.Must(t, got.Close()) }()

	if needSnapshot {
		noerror.Must(t, noerror.NotEqual(nil).Actual(expected.Data))

		var actual interface{}
		if expected.ModifyResponse != nil {
			actual = expected.ModifyResponse(got)
		} else {
			actual = got.JSONData()
		}

		noerror.Should(t,
			jsonequal.ShouldBeSame(
				jsonequal.From(actual),
				jsonequal.From(expected.Data),
			),
		)
	}
}
