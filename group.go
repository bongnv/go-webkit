package gwf

import "github.com/julienschmidt/httprouter"

// Group is a group of sub-routes. It can be used for a group of routes which share same middlewares.
type Group struct {
	prefix       string
	routeOptions []RouteOption
	router       *httprouter.Router
	logger       Logger
}

// GET registers a new GET route for a path with handler.
func (g *Group) GET(path string, h Handler, opts ...RouteOption) {
	r := g.newRoute(h, opts)
	g.router.GET(g.buildPath(path), r.buildHandle())
}

// POST registers a new POST route for a path with handler.
func (g *Group) POST(path string, h Handler, opts ...RouteOption) {
	r := g.newRoute(h, opts)
	g.router.POST(g.buildPath(path), r.buildHandle())
}

// PUT registers a new PUT route for a path with handler.
func (g *Group) PUT(path string, h Handler, opts ...RouteOption) {
	r := g.newRoute(h, opts)
	g.router.PUT(g.buildPath(path), r.buildHandle())
}

// PATCH registers a new PATCH route for a path with handler.
func (g *Group) PATCH(path string, h Handler, opts ...RouteOption) {
	r := g.newRoute(h, opts)
	g.router.PATCH(g.buildPath(path), r.buildHandle())

}

// DELETE registers a new DELETE route for a path with handler.
func (g *Group) DELETE(path string, h Handler, opts ...RouteOption) {
	r := g.newRoute(h, opts)
	g.router.DELETE(g.buildPath(path), r.buildHandle())
}

// newRoute creates a new route give Handler and a list of RouteOption.
func (g *Group) newRoute(h Handler, opts []RouteOption) *route {
	r := &route{
		errorHandler: defaultErrorHandler(g.logger),
		handler:      h,
		logger:       g.logger,
		transformers: []handleTransformer{
			brwTransformer,
		},
	}

	r.applyOpts(g.routeOptions)
	r.applyOpts(opts)

	return r
}

func (g *Group) buildPath(path string) string {
	return g.prefix + path
}
