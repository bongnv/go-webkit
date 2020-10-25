package webkit

// Logger defines a Logger.
type Logger interface {
	// Println prints out logs like fmt.Println.
	Println(...interface{})
}
