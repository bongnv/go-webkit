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

	"github.com/julienschmidt/httprouter"
)

// New creates a new application.
func New() *Application {
	return &Application{
		router:     httprouter.New(),
		port:       8080,
		shutdownCh: make(chan struct{}),
	}
}

// Application is a web application.
type Application struct {
	port int

	router       *httprouter.Router
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
	app.router.GET(path, buildHandlerFunc(h))
}

// POST registers a new POST route for a path with handler.
func (app *Application) POST(path string, h Handler) {
	app.router.POST(path, buildHandlerFunc(h))
}

// PUT registers a new PUT route for a path with handler.
func (app *Application) PUT(path string, h Handler) {
	app.router.PUT(path, buildHandlerFunc(h))
}

// PATCH registers a new PATCH route for a path with handler.
func (app *Application) PATCH(path string, h Handler) {
	app.router.PATCH(path, buildHandlerFunc(h))

}

// DELETE registers a new DELETE route for a path with handler.
func (app *Application) DELETE(path string, h Handler) {
	app.router.DELETE(path, buildHandlerFunc(h))
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
