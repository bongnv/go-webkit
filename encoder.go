package gwf

import (
	"encoding/json"
	"net/http"
)

// Encoder define a request decoder.
type Encoder interface {
	// Encode encodes obj and writes to http.ResponseWriter.
	Encode(w http.ResponseWriter, obj interface{}) error
}

// WithEncoder specifies the encoder which will be used to encode payload to HTTP response.
func WithEncoder(e Encoder) RouteOptionFn {
	return func(r *route) {
		r.encoder = e
	}
}

type defaultEncoder struct{}

func (d *defaultEncoder) Encode(w http.ResponseWriter, obj interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(obj)
}

func newEncoder() Encoder {
	return &defaultEncoder{}
}
