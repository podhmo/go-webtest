package replace

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonpointer"
)

// ErrUnsupportedType :
var ErrUnsupportedType = fmt.Errorf("unsupported type")

func coerce(v interface{}) (interface{}, error) {
	switch v := v.(type) {
	case map[string]interface{}:
		return v, nil
	case []interface{}:
		return v, nil
	case *[]interface{}:
		if v == nil {
			return []interface{}{}, nil
		}
		return *v, nil
	case *map[string]interface{}:
		if v == nil {
			return map[string]interface{}{}, nil
		}
		return *v, nil
	case *interface{}:
		if v == nil {
			return map[string]interface{}{}, nil // xxx:
		}
		return *v, nil
	default:
		return nil, errors.WithMessagef(ErrUnsupportedType, "only map[string]interface{} and []interface{}. this is %T", v)
	}
}

// ByMap replace data by map
func ByMap(data interface{}, refMap map[string]interface{}) (interface{}, error) {
	data, err := coerce(data)
	if err != nil {
		return nil, err
	}
	for k, v := range refMap {
		jptr, err := gojsonpointer.NewJsonPointer(strings.TrimPrefix(k, "#"))
		if err != nil {
			return nil, errors.WithMessagef(err, "parse %q as jsonpointer", k)
		}
		if _, err := jptr.Set(data, v); err != nil {
			return nil, errors.WithMessagef(err, "access %q on data (set)", k)
		}
	}
	return data, nil
}

// ByPalette replace data by reference array and palette
func ByPalette(data interface{}, refs []string, palette interface{}) (interface{}, error) {
	data, err := coerce(data)
	if err != nil {
		return nil, err
	}
	for _, k := range refs {
		jptr, err := gojsonpointer.NewJsonPointer(strings.TrimPrefix(k, "#"))
		if err != nil {
			return nil, errors.WithMessagef(err, "parse %q as jsonpointer", k)
		}
		v, _, err := jptr.Get(palette)
		if err != nil {
			return nil, errors.WithMessagef(err, "access %q on pallete (get)", k)
		}
		if _, err := jptr.Set(data, v); err != nil {
			return nil, errors.WithMessagef(err, "access %q on data (set)", k)
		}
	}
	return data, nil
}
