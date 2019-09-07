package jsonequal

import (
	"fmt"

	"github.com/nsf/jsondiff"
)

// FailPlain :
func FailPlain(
	left interface{},
	right interface{},
	lb []byte,
	rb []byte,
) error {
	ls, rs := string(lb), string(rb)
	if ls == rs {
		msg := "not equal json\nleft (%[1]T):\n	%[3]s\nright (%[2]T):\n	%[4]s"
		return fmt.Errorf(msg, left, right, ls, rs)
	}
	// todo : more redable expression
	msg := "not equal json\nleft:\n	%[3]s\nright:\n	%[4]s"
	return fmt.Errorf(msg, left, right, ls, rs)
}

// FailJSONDiff :
func FailJSONDiff(
	left interface{},
	right interface{},
	lb []byte,
	rb []byte,
) error {
	options := jsondiff.DefaultConsoleOptions()
	diff, s := jsondiff.Compare(lb, rb, &options)
	return fmt.Errorf("%s\n%s\n%s", diff.String(), s, FailPlain(left, right, lb, rb).Error())
}
