package testclient

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type withPrefixWriter struct {
	Prefix string
	Writer io.Writer
}

func (w *withPrefixWriter) Write(b []byte) (int, error) {
	buf := bytes.NewBuffer(b)
	s := bufio.NewScanner(buf)

	var total int
	for s.Scan() {
		n, err := io.WriteString(w.Writer, w.Prefix)
		if err != nil {
			return total + n, err
		}
		m, err := w.Writer.Write(s.Bytes())
		if err != nil {
			return total + n + m, err
		}
		_, err = io.WriteString(w.Writer, "\n")
		if err != nil {
			return total + n + m + 1, err
		}
		total += n + m + 1
	}
	return total, nil
}

// DebugRoundTripper :
type DebugRoundTripper struct {
	IgnoreDumpRequest  bool
	IgnoreDumpResponse bool
	Quiet              bool

	Writer    io.Writer
	Transport http.RoundTripper
}

// Decorate :
func (d *DebugRoundTripper) Decorate(transport http.RoundTripper) RoundTripperDecorator {
	new := *d
	if new.Transport != nil {
		log.Printf("!! %T.Transport is not nil, overwrite original one", d)
	}
	new.Transport = transport
	return &new
}

// transport :
func (d *DebugRoundTripper) transport() http.RoundTripper {
	if d.Transport != nil {
		return d.Transport
	}
	return http.DefaultTransport
}

// writer :
func (d *DebugRoundTripper) writer() io.Writer {
	if d.Writer != nil {
		return d.Writer
	}
	return &withPrefixWriter{Writer: os.Stderr, Prefix: "\t"}
}

// RoundTrip :
func (t *DebugRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if !t.IgnoreDumpRequest {
		b, err := httputil.DumpRequest(req, !t.Quiet)
		if err != nil {
			return nil, err
		}
		w := t.writer()
		fmt.Fprintln(w, "\x1b[34mRequest : ------------------------------\x1b[0m")
		if _, err := w.Write(b); err != nil {
			return nil, err
		}
		fmt.Fprintln(w, "\x1b[34m----------------------------------------\x1b[0m")
	}

	resp, err = t.transport().RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if !t.IgnoreDumpResponse {
		b, err := httputil.DumpResponse(resp, !t.Quiet)
		if err != nil {
			return nil, err
		}
		w := t.writer()
		fmt.Fprintln(w, "\x1b[32mResponse: ------------------------------\x1b[0m")
		if _, err := w.Write(b); err != nil {
			return nil, err
		}
		fmt.Fprintln(w, "\x1b[32m----------------------------------------\x1b[0m")
	}
	return resp, nil
}
