package nanny

// Option defines an application Option.
type Option interface {
	Apply(app *Application)
}

// OptionFn defines a function that implements Option
type OptionFn func(app *Application)

// Apply implements Option.
func (opt OptionFn) Apply(app *Application) {
	opt(app)
}

// WithAddress specifies the TCP address for the server to listen on,
func WithAddress(addr string) OptionFn {
	return func(app *Application) {
		app.addr = addr
	}
}
