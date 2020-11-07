package gwf

// Plugin is a set of Option to enrich an application.
type Plugin []Option

// Apply implements Option interface. It applies options to the application.
func (p Plugin) Apply(app *Application) {
	app.applyOpts(p)
}

// DefaultApp is a plugin to provide a set of common options for an application.
var DefaultApp Plugin = []Option{
	WithLogger(defaultLogger()),
	WithRecovery(),
	WithCORS(DefaultCORSConfig),
	WithGzip(DefaultGzipConfig),
}
