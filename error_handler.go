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
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(errResp.Error())); err != nil {
			logger.Println("Error", err, "while encoding", errResp)
		}
	}
}
