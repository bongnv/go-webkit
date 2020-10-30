package webkit

import (
	"log"
	"os"
)

// Logger defines a Logger.
type Logger interface {
	// Println prints out logs like fmt.Println.
	Println(...interface{})
}

func defaultLogger() Logger {
	return log.New(os.Stderr, "", log.LstdFlags)
}
