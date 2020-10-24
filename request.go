package webkit

import "net/http"

// Request defines a HTTP request.
type Request interface {}

type requestImpl struct {
  httpWriter http.ResponseWriter
  httpReq *http.Request
}

func (r *requestImpl) responseError(err error) {
  r.httpWriter.WriteHeader(http.StatusInternalServerError)
  // TODO: Add logs here
  _, _ = r.httpWriter.Write([]byte(err.Error()))
}
