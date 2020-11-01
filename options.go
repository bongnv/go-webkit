package gwf

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

// WithLogger specifies a custom Logger for tha application.
func WithLogger(l Logger) OptionFn {
	return func(app *Application) {
		if l != nil {
			app.logger = l
			app.MustRegister("logger", l)
		}
	}
}
