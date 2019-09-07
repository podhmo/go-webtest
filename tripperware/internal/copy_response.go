package internal

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// CopyResponse has Side effect
func CopyResponse(res *http.Response) *http.Response {
	if res == nil {
		return nil
	}
	copied := *res
	if res.Body == http.NoBody {
		return &copied
	}
	var b bytes.Buffer
	res.Body = ioutil.NopCloser(io.TeeReader(res.Body, &b))
	copied.Body = ioutil.NopCloser(&b)
	return &copied
}
