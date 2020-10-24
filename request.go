package webkit

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Request defines a HTTP request.
type Request interface {
	// Response sends take resp and render it to clients.
	Response(resp interface{}) error
}

type requestImpl struct {
	httpWriter http.ResponseWriter
	httpReq    *http.Request
	params     httprouter.Params
}

func (r *requestImpl) responseError(err error) {
	r.httpWriter.WriteHeader(http.StatusInternalServerError)
	// TODO: Add logs here
	_, _ = r.httpWriter.Write([]byte(err.Error()))
}

func (r *requestImpl) Response(resp interface{}) error {
	enc := json.NewEncoder(r.httpWriter)
	return enc.Encode(resp)
}
