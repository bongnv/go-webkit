package webkit

import (
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
func (fn *RouteOptionFn) Apply(app *Application) {
	app.routeOptions = append(app.routeOptions, fn)
}

type route struct {
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
		req := &requestImpl{
			httpWriter: w,
			httpReq:    httpReq,
		}

		err := h(ctx, req)
		if err != nil {
			req.responseError(err)
		}
	}
}
