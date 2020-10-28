package webkit

import (
	"context"
	"net/http"

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

type route struct {
	decoder     Decoder
	encoder     Encoder
	handler     Handler
	logger      Logger
	middlewares []Middleware
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

	return func(w http.ResponseWriter, httpReq *http.Request, params httprouter.Params) {
		ctx := httpReq.Context()
		ctx = context.WithValue(ctx, ctxKeyHTTPResponseWriter, w)

		req := &requestImpl{
			decoder: r.decoder,
			httpReq: httpReq,
			params:  params,
		}

		resp, err := h(ctx, req)
		if err == nil {
			err = r.writeToHTTPResponse(w, resp)
		}

		if err != nil {
			// TODO: Add error handler
			r.responseError(w, err)
		}
	}
}

func (r *route) writeToHTTPResponse(w http.ResponseWriter, resp interface{}) error {
	if resp == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	return r.encoder.Encode(w, resp)
}

func (r *route) responseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	// TODO: Add logs here
	_, _ = w.Write([]byte(err.Error()))
}
