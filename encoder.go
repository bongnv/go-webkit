package webkit

import (
	"encoding/json"
	"net/http"
)

// Encoder define a request decoder.
type Encoder interface {
	// Encode encodes obj and writes to http.ResponseWriter.
	Encode(w http.ResponseWriter, obj interface{}) error
}

type defaultEncoder struct{}

func (d *defaultEncoder) Encode(w http.ResponseWriter, obj interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(obj)
}

func newEncoder() Encoder {
	return &defaultEncoder{}
}
