package nanny

import (
	"encoding/json"
	"net/http"
)

const (
	jsonScheme = "application/json"
)

// Encoder define a request decoder.
type Encoder interface {
	// Encode encodes obj and writes to http.ResponseWriter.
	Encode(w http.ResponseWriter, resp interface{}) error
}

// WithEncoder specifies the encoder which will be used to encode payload to HTTP response.
func WithEncoder(e Encoder) RouteOptionFn {
	return func(r *route) {
		r.encoder = e
	}
}

type defaultEncoder struct{}

func (e defaultEncoder) Encode(w http.ResponseWriter, resp interface{}) error {
	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	if customResp, ok := resp.(CustomHTTPResponse); ok {
		customResp.WriteTo(w)
		return nil
	}

	w.Header().Add(HeaderContentType, jsonScheme)
	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}
