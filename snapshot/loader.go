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

// Extra is extra data
type Extra struct {
	ReplaceMap map[string]interface{} // data replacement setting, on loading
	Metadata   map[string]interface{}
}

// Loader :
type Loader struct {
	Encode func(io.Writer, interface{}, *Extra) error
	Decode func(io.Reader, interface{}, *Extra) error
}

// Save :
func (r *Loader) Save(fpath string, val interface{}, extra *Extra) (err error) {
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
	return r.Encode(wf, val, extra)
}

func (r *Loader) Load(fpath string, want interface{}, extra *Extra) (err error) {
	rf, err := os.Open(fpath)
	if err != nil {
		return errors.WithMessage(err, "on open file")
	}
	defer func() {
		if cerr := rf.Close(); cerr != nil {
			err = cerr
		}
	}()
	if err := r.Decode(rf, want, extra); err != nil {
		return errors.WithMessage(err, "on decoder.Decode")
	}
	return nil
}

// todo: add mtime?

type saveData struct {
	ModifiedAt time.Time              `json:"modifiedAt"`
	Replaced   map[string]interface{} `json:"replaced,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Data       interface{}            `json:"data"`
}

type loadData struct {
	ModifiedAt time.Time              `json:"modifiedAt"`
	Replaced   map[string]interface{} `json:"replaced,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Data       json.RawMessage        `json:"data"`
}

// NewJSONLoader :
func NewJSONLoader() *Loader {
	return &Loader{
		Encode: func(w io.Writer, val interface{}, extra *Extra) error {
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)
			data := &saveData{
				ModifiedAt: time.Now(),
				Data:       val,
				Replaced:   extra.ReplaceMap,
				Metadata:   extra.Metadata,
			}
			if err := encoder.Encode(data); err != nil {
				return errors.WithMessage(err, "on json encode")
			}
			return nil
		},
		Decode: func(r io.Reader, val interface{}, extra *Extra) error {
			decoder := json.NewDecoder(r)
			var data loadData
			if err := decoder.Decode(&data); err != nil {
				return errors.WithMessage(err, "on json decode")
			}
			if err := json.Unmarshal(data.Data, val); err != nil {
				return errors.WithMessage(err, "on unmarshal raw message")
			}
			if extra.ReplaceMap == nil {
				return nil
			}
			_, err := replace.ByMap(val, extra.ReplaceMap)
			if err != nil {
				return errors.WithMessage(err, "on replace data by map")
			}
			return err
		},
	}
}
