package nanny

import "net/http"

// ErrorHandler defines a handler which handles error.
type ErrorHandler func(w http.ResponseWriter, err error) error

// WithErrorHandler is a RouteOption to specify a custom ErrorHandler.
func WithErrorHandler(errHandler ErrorHandler) RouteOptionFn {
	return func(r *route) {
		if errHandler != nil {
			r.errorHandler = errHandler
		}
	}
}

func defaultErrorHandler() ErrorHandler {
	return func(w http.ResponseWriter, errResp error) error {
		code := http.StatusInternalServerError
		body := []byte(errResp.Error())

		if customResp, ok := errResp.(CustomHTTPResponse); ok {
			customResp.WriteTo(w)
			return nil
		}

		w.WriteHeader(code)
		_, err := w.Write(body)
		return err
	}
}
