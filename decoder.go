package nanny

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
)

// Decoder defines a request decoder.
type Decoder interface {
	// Decode decodes a request to a struct. req.PostForm is called in advanced.
	Decode(obj interface{}, req *http.Request) error
}

// WithDecoder specifies the decoder which will be used.
func WithDecoder(d Decoder) RouteOptionFn {
	return func(r *route) {
		r.decoder = d
	}
}

func newDecoder() Decoder {
	return &defaultDecoder{
		schemaDecoder: schema.NewDecoder(),
	}
}

type defaultDecoder struct {
	schemaDecoder *schema.Decoder
}

func (d *defaultDecoder) Decode(obj interface{}, req *http.Request) error {
	if err := d.schemaDecoder.Decode(obj, req.Form); err != nil {
		return err
	}

	contentType := req.Header.Get("Content-Type")
	if contentType == "application/json" && req.ContentLength > 0 {
		// erase all values in forms so that they won't overwrite parsed json values
		jsonDecoder := json.NewDecoder(req.Body)
		if err := jsonDecoder.Decode(obj); err != nil {
			return err
		}
	}

	return nil
}
