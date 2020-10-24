package webkit

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// New creates a new application.
func New() *Application {
	return &Application{
		router: mux.NewRouter(),
		port:   8080,
	}
}

// Application is a web application.
type Application struct {
	router *mux.Router
	port   int
}

// Run starts an HTTP server.
func (app *Application) Run() error {
	srv := &http.Server{
		Handler: app.router,
		Addr:    fmt.Sprint(":", app.port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Serving at port", app.port)
	log.Fatal(srv.ListenAndServe())
	return nil
}

// GET registers a new GET route for a path with handler.
func (app *Application) GET(path string, h Handler) {
	route := app.router.Path(path).Methods(http.MethodGet)
	route.HandlerFunc(buildHandlerFunc(h))
}
