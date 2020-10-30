package webkit

import "encoding/json"

// HTTPError is a simple implementation of HTTP Error.
type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

// HTTPResponse returns HTTP code and body for CustomHTTPResponse.
func (err *HTTPError) HTTPResponse() (int, []byte) {
	data, errMarshal := json.Marshal(err)
	if errMarshal != nil {
		panic(errMarshal)
	}

	return err.Code, data
}

// Error implements error interface.
func (err *HTTPError) Error() string {
	return err.Message
}
