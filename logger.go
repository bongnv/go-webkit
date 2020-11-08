package nanny

import (
	"log"
	"os"
)

// Logger defines a Logger.
type Logger interface {
	// Println prints out logs like fmt.Println.
	Println(...interface{})
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

func defaultLogger() Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}
