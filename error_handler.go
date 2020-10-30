package webkit

import "net/http"

// ErrorHandler defines a handler which handles error.
type ErrorHandler func(w http.ResponseWriter, err error)

// WithErrorHandler is a RouteOption to specify a custom ErrorHandler.
func WithErrorHandler(errHandler ErrorHandler) RouteOptionFn {
	return func(r *route) {
		if errHandler != nil {
			r.errorHandler = errHandler
		}
	}
}

func defaultErrorHandler(logger Logger) ErrorHandler {
	return func(w http.ResponseWriter, errResp error) {
		code := http.StatusInternalServerError
		body := []byte(errResp.Error())

		if withCustomResp, ok := errResp.(CustomHTTPResponse); ok {
			code, body = withCustomResp.HTTPResponse()
		}

		w.WriteHeader(code)
		if _, err := w.Write(body); err != nil {
			logger.Println("Error", err, "write sending response", err)
		}
	}
}
