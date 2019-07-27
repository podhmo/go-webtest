package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Recorder :
type Recorder struct {
	Exists func(string) bool
	Path   func(testing.TB) string
	Loader *Loader
}

func NewTestdataRecorder(loader *Loader) *Recorder {
	return &Recorder{
		Exists: func(fpath string) bool {
			_, err := os.Stat(fpath)
			return err == nil
		},
		Path: func(t testing.TB) string {
			return fmt.Sprintf("%s.golden", filepath.Join("testdata", t.Name()))
		},
		Loader: loader,
	}
}

// Config :
type Config struct {
	Recorder  *Recorder
	Overwrite bool
	self      string
}

// Run :
func (c *Config) Run(
	t testing.TB,
	got interface{},
) interface{} {
	r := c.Recorder
	fpath := r.Path(t)
	existed := r.Exists(fpath)
	if !existed || c.Overwrite || c.self == fpath || c.self == filepath.Base(fpath) {
		t.Logf("save testdata: %q", fpath)
		if err := r.Loader.Save(fpath, got); err != nil {
			t.Fatalf("record: %s", err)
		}
	}
	t.Logf("load testdata: %q", fpath)

	var want interface{}
	if err := r.Loader.Load(fpath, &want); err != nil {
		t.Fatalf("replay: %s", err)
	}
	return want
}

// default recorder
var recorder = NewTestdataRecorder(NewJSONLoader())

// Take snapshot tests if needed and return expected data
func Take(
	t testing.TB,
	got interface{},
	options ...func(*Config),
) interface{} {
	c := &Config{
		Recorder:  recorder,
		Overwrite: false,
	}

	// default overwrite
	WithUpdateByEnvvar("SNAPSHOT")(c)

	for _, opt := range options {
		opt(c)
	}
	return c.Run(t, got)
}

// WithForceUpdate :
func WithForceUpdate() func(*Config) {
	return func(c *Config) {
		c.Overwrite = true
	}
}

// WithUpdateByEnvvar :
func WithUpdateByEnvvar(s string) func(*Config) {
	return func(c *Config) {
		v := strings.Trim(os.Getenv(s), " ")
		if v == "" {
			return
		}
		if v == "1" {
			c.Overwrite = true
			return
		}
		c.self = v
	}
}