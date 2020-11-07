package nanny

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type contextKey int

const (
	_ contextKey = iota
	ctxKeyHTTPResponseWriter
	ctxKeyApp
	ctxKeyRoute
)

// ResponseHeaderFromCtx returns Header for HTTP response which will be sent.
// The function returns a nil map if the Header doesn't exist.
func ResponseHeaderFromCtx(ctx context.Context) http.Header {
	w, ok := ctx.Value(ctxKeyHTTPResponseWriter).(http.ResponseWriter)
	if ok {
		return w.Header()
	}

	return nil
}

func contextInjector() OptionFn {
	return func(app *Application) {
		var routeOpt RouteOptionFn = func(r *route) {
			routeTransformer := func(next httprouter.Handle) httprouter.Handle {
				return func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
					ctx := req.Context()
					ctx = context.WithValue(ctx, ctxKeyApp, app)
					ctx = context.WithValue(ctx, ctxKeyRoute, r)
					ctx = context.WithValue(ctx, ctxKeyHTTPResponseWriter, rw)
					req = req.WithContext(ctx)

					next(rw, req, p)
				}
			}

			r.transformers = append(r.transformers, routeTransformer)
		}

		app.routeOptions = append(app.routeOptions, routeOpt)
	}
}

func loggerFromCtx(ctx context.Context) Logger {
	app, ok := ctx.Value(ctxKeyApp).(*Application)
	if !ok {
		return nil
	}

	return app.logger
}
