package nanny

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// RouteOption defines an option to customize a route.
type RouteOption interface {
	ApplyRoute(r *route)
}

// RouteOptionFn defines a function implementation of RouteOption.
type RouteOptionFn func(r *route)

// ApplyRoute implements RouteOption.
func (fn RouteOptionFn) ApplyRoute(r *route) {
	fn(r)
}

// Apply implements Option.
func (fn RouteOptionFn) Apply(app *Application) {
	app.routeOptions = append(app.routeOptions, fn)
}

// CustomHTTPResponse defines an interface to support custom HTTP response.
type CustomHTTPResponse interface {
	WriteTo(w http.ResponseWriter)
}

type handleTransformer func(httprouter.Handle) httprouter.Handle

type route struct {
	decoder      Decoder
	encoder      Encoder
	errorHandler ErrorHandler
	handler      Handler
	logger       Logger
	method       string
	middlewares  []Middleware
	path         string
	timeout      time.Duration
	transformers []handleTransformer
}

func (r *route) applyOpts(opts []RouteOption) {
	for _, o := range opts {
		o.ApplyRoute(r)
	}
}

func (r *route) buildHandle() httprouter.Handle {
	h := r.handler
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}

	handle := func(w http.ResponseWriter, httpReq *http.Request, params httprouter.Params) {
		ctx := httpReq.Context()

		req := &requestImpl{
			decoder: r.decoder,
			httpReq: httpReq,
			params:  params,
		}

		resp, err := h(ctx, req)
		if err != nil {
			if errHandle := r.errorHandler(w, err); errHandle != nil {
				r.logger.Println("Error", errHandle, "while handling error")
			}
			return
		}

		if errWrite := r.encoder.Encode(w, resp); errWrite != nil {
			r.logger.Println("Error", errWrite, "while sending response")
		}
	}

	for i := len(r.transformers) - 1; i >= 0; i-- {
		handle = r.transformers[i](handle)
	}

	return handle
}
