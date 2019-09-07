package jsonequal

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/yudai/gojsondiff"
	"github.com/yudai/gojsondiff/formatter"
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
	err := failJSONDiff(left, right, lb, rb)
	if err == nil {
		return FailPlain(left, right, lb, rb)
	}
	if os.Getenv("DEBUG") != "" {
		return errors.WithMessage(err, FailPlain(left, right, lb, rb).Error())
	}
	return err
}

// failJSONDiff :
func failJSONDiff(
	left interface{},
	right interface{},
	lb []byte,
	rb []byte,
) (err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = errors.WithMessage(FailPlain(left, right, lb, rb), fmt.Sprintf("%s", rerr))
		}
	}()

	config := formatter.AsciiFormatterConfig{
		ShowArrayIndex: true,
		Coloring:       true,
	}

	switch left := left.(type) {
	case map[string]interface{}:
		right, ok := right.(map[string]interface{})
		if !ok {
			return nil
		}
		d := gojsondiff.New()
		diff := d.CompareObjects(left, right)
		formatter := formatter.NewAsciiFormatter(left, config)
		s, err := formatter.Format(diff)
		if err != nil {
			err = errors.WithMessage(fmt.Errorf("diff%s", s), err.Error())
		}
		return fmt.Errorf("diff%s", s)
	case []interface{}:
		right, ok := right.([]interface{})
		if !ok {
			return nil
		}
		d := gojsondiff.New()
		diff := d.CompareArrays(left, right)
		formatter := formatter.NewAsciiFormatter(left, config)
		s, err := formatter.Format(diff)
		if err != nil {
			err = errors.WithMessage(fmt.Errorf("diff%s", s), err.Error())
		}
		return fmt.Errorf("diff%s", s)
	default:
		return nil
	}
}
