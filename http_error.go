package nanny

import (
	"encoding/json"
	"net/http"
)

// HTTPError is a simple implementation of HTTP Error.
type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

// WriteTo implements CustomHTTPResponse. It encodes the response as JSON format.
func (err *HTTPError) WriteTo(w http.ResponseWriter) {
	w.Header().Add(HeaderContentType, jsonScheme)
	w.WriteHeader(err.Code)
	enc := json.NewEncoder(w)
	_ = enc.Encode(err)
}

// Error implements error interface.
func (err *HTTPError) Error() string {
	return err.Message
}
