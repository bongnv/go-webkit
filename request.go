package webkit

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Request defines a HTTP request.
type Request interface {
	// Decode decodes the request to an object.
	Decode(obj interface{}) error
	// Response sends take resp and render it to clients.
	Response(resp interface{}) error
}

type requestImpl struct {
	decoder    Decoder
	httpWriter http.ResponseWriter
	httpReq    *http.Request
	params     httprouter.Params
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

func (r *requestImpl) Response(resp interface{}) error {
	enc := json.NewEncoder(r.httpWriter)
	return enc.Encode(resp)
}

func (r *requestImpl) responseError(err error) {
	r.httpWriter.WriteHeader(http.StatusInternalServerError)
	// TODO: Add logs here
	_, _ = r.httpWriter.Write([]byte(err.Error()))
}
