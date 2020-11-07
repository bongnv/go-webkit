package gwf

import (
	"net/http"
)

// RouteGroup is a group of sub-routes. It can be used for a group of routes which share same middlewares.
type RouteGroup struct {
	prefix       string
	routeOptions []RouteOption

	app *Application
}

// GET registers a new GET route for a path with handler.
func (g *RouteGroup) GET(path string, h Handler, opts ...RouteOption) {
	g.addRoute(http.MethodGet, path, h, opts)
}

// POST registers a new POST route for a path with handler.
func (g *RouteGroup) POST(path string, h Handler, opts ...RouteOption) {
	g.addRoute(http.MethodPost, path, h, opts)

}

// PUT registers a new PUT route for a path with handler.
func (g *RouteGroup) PUT(path string, h Handler, opts ...RouteOption) {
	g.addRoute(http.MethodPut, path, h, opts)
}

// PATCH registers a new PATCH route for a path with handler.
func (g *RouteGroup) PATCH(path string, h Handler, opts ...RouteOption) {
	g.addRoute(http.MethodPatch, path, h, opts)
}

// DELETE registers a new DELETE route for a path with handler.
func (g *RouteGroup) DELETE(path string, h Handler, opts ...RouteOption) {
	g.addRoute(http.MethodDelete, path, h, opts)
}

// Group creates a group of sub-routes
func (g *RouteGroup) Group(prefix string, opts ...RouteOption) *RouteGroup {
	if len(prefix) == 0 || prefix[0] != '/' {
		panic("path must begin with '/' in path '" + prefix + "'")
	}

	// Strip trailing / (if present) as all added sub paths must start with a /
	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	return &RouteGroup{
		prefix:       g.prefix + prefix,
		routeOptions: append(g.routeOptions, opts...),
		app:          g.app,
	}
}

func (g *RouteGroup) addRoute(method, path string, h Handler, opts []RouteOption) {
	r := &route{
		errorHandler: defaultErrorHandler(g.app.logger),
		handler:      h,
		logger:       g.app.logger,
		transformers: []handleTransformer{
			brwTransformer,
		},
		method: method,
		path:   g.prefix + path,
	}

	r.applyOpts(g.routeOptions)
	r.applyOpts(opts)
	g.app.routes = append(g.app.routes, r)
}
