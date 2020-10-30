package gwf

// Middleware defines a middleware to provide additional logic.
type Middleware func(Handler) Handler

// ApplyRoute implements RouteOption.
func (m Middleware) ApplyRoute(r *route) {
	r.middlewares = append(r.middlewares, m)
}

// Apply implements Option.
func (m Middleware) Apply(app *Application) {
	app.routeOptions = append(app.routeOptions, m)
}
