package snapshot

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/podhmo/go-webtest/replace"
)

// todo: crete paramameter object? (include repMap)

// Loader :
type Loader struct {
	Encode func(io.Writer, interface{}, map[string]interface{}) error
	Decode func(io.Reader, interface{}, map[string]interface{}) error
}

// Save :
func (r *Loader) Save(fpath string, val interface{}, repMap map[string]interface{}) (err error) {
	if err := os.Mkdir(filepath.Dir(fpath), 0744); err != nil {
		if !os.IsExist(err) {
			return errors.WithMessagef(err, "create testdata directory, %q", filepath.Dir(fpath))
		}
	}
	wf, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := wf.Close(); cerr != nil {
			err = cerr
		}
	}()
	return r.Encode(wf, val, repMap)
}

func (r *Loader) Load(fpath string, want interface{}, repMap map[string]interface{}) (err error) {
	rf, err := os.Open(fpath)
	if err != nil {
		return errors.WithMessage(err, "on open file")
	}
	defer func() {
		if cerr := rf.Close(); cerr != nil {
			err = cerr
		}
	}()
	if err := r.Decode(rf, want, repMap); err != nil {
		return errors.WithMessage(err, "on decoder.Decode")
	}
	return nil
}

// todo: add mtime?

type saveData struct {
	ModifiedAt time.Time              `json:"modifiedAt"`
	Replaced   map[string]interface{} `json:"replaced,omitempty"`
	Data       interface{}            `json:"data"`
}

type loadData struct {
	ModifiedAt time.Time              `json:"modifiedAt"`
	Replaced   map[string]interface{} `json:"replaced,omitempty"`
	Data       json.RawMessage        `json:"data"`
}

// NewJSONLoader :
func NewJSONLoader() *Loader {
	return &Loader{
		Encode: func(w io.Writer, val interface{}, repMap map[string]interface{}) error {
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)
			data := &saveData{
				ModifiedAt: time.Now(),
				Data:       val,
				Replaced:   repMap,
			}
			if err := encoder.Encode(data); err != nil {
				return errors.WithMessage(err, "on json encode")
			}
			return nil
		},
		Decode: func(r io.Reader, val interface{}, repMap map[string]interface{}) error {
			decoder := json.NewDecoder(r)
			var data loadData
			if err := decoder.Decode(&data); err != nil {
				return errors.WithMessage(err, "on json decode")
			}
			if err := json.Unmarshal(data.Data, val); err != nil {
				return errors.WithMessage(err, "on unmarshal raw message")
			}
			if repMap == nil {
				return nil
			}
			_, err := replace.ByMap(val, repMap)
			if err != nil {
				return errors.WithMessage(err, "on replace data by map")
			}
			return err
		},
	}
}
