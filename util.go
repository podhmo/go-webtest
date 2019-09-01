package webtest

import "net/url"

// MustParseQuery :
func MustParseQuery(query string) url.Values {
	vals, err := url.ParseQuery(query)
	if err != nil {
		panic(err)
	}
	return vals
}
