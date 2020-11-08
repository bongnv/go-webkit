package nanny

// Plugin is a set of Option to enrich an application.
type Plugin []Option

// Apply implements Option interface. It applies options to the application.
func (p Plugin) Apply(app *Application) {
	app.applyOpts(p)
}
