package webkit

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
	// ResponseHeader returns the header map that will be sent.
	ResponseHeader() http.Header
}

type requestImpl struct {
	decoder    Decoder
	httpWriter http.ResponseWriter
	httpReq    *http.Request
	params     httprouter.Params
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

func (r *requestImpl) ResponseHeader() http.Header {
	return r.httpWriter.Header()
}

func (r *requestImpl) responseError(err error) {
	r.httpWriter.WriteHeader(http.StatusInternalServerError)
	// TODO: Add logs here
	_, _ = r.httpWriter.Write([]byte(err.Error()))
}
