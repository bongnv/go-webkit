package webkit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// New creates a new application.
func New() *Application {
	return &Application{
		router:     mux.NewRouter(),
		port:       8080,
		shutdownCh: make(chan struct{}),
	}
}

// Application is a web application.
type Application struct {
	port int

	router       *mux.Router
	shutdownCh   chan struct{}
	shutdownOnce sync.Once
	srv          *http.Server
	wg           sync.WaitGroup
}

// Run starts an HTTP server.
func (app *Application) Run() error {
	app.srv = &http.Server{
		Handler: app.router,
		Addr:    fmt.Sprint(":", app.port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	app.execute(func() {
		log.Println("Serving at port", app.port)
		if err := app.srv.ListenAndServe(); err != nil {
			log.Println("Error:", err)
			app.shutdown()
		}
	})

	app.setupGracefulShutdown()

	app.wg.Wait()
	return nil
}

// GET registers a new GET route for a path with handler.
func (app *Application) GET(path string, h Handler) {
	route := app.router.Path(path).Methods(http.MethodGet)
	route.HandlerFunc(buildHandlerFunc(h))
}

// execute starts a function in a goroutine.
// It tracks the execution in a WaitGroup for graceful shutdown.
func (app *Application) execute(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		fn()
	}()
}

// setupGracefulShutdown starts a goroutine for interrupt signal and proceed with graceful shutdown.
func (app *Application) setupGracefulShutdown() {
	app.execute(func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		select {
		case <-sigint:
			app.shutdown()
		case <-app.shutdownCh:
		}
	})
}

func (app *Application) shutdown() {
	app.shutdownOnce.Do(func() {
		// We received an interrupt signal, shut down.
		if err := app.srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(app.shutdownCh)
	})
}
