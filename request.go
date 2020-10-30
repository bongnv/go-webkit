package gwf

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Request defines a HTTP request.
type Request interface {
	// HTTPRequest returns the http.Request.
	HTTPRequest() *http.Request
	// Decode decodes the request to an object.
	Decode(obj interface{}) error
}

type requestImpl struct {
	decoder Decoder
	httpReq *http.Request
	params  httprouter.Params
}

func (r *requestImpl) HTTPRequest() *http.Request {
	return r.httpReq
}

func (r *requestImpl) Decode(obj interface{}) error {
	if err := r.httpReq.ParseForm(); err != nil {
		return err
	}

	for _, p := range r.params {
		r.httpReq.Form.Set(p.Key, p.Value)
	}

	return r.decoder.Decode(obj, r.httpReq)
}
