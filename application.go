package webkit

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Handler defines a function to serve HTTP requests.
type Handler func(ctx context.Context, req Request) error

// New creates a new application.
func New(opts ...Option) *Application {
	app := &Application{
		decoder:       newDecoder(),
		router:        httprouter.New(),
		port:          8080,
		readyCh:       make(chan struct{}),
		srvShutdownCh: make(chan struct{}),
		logger:        log.New(os.Stderr, "", log.LstdFlags),
	}

	app.applyOpts(opts)
	return app
}

// Default returns an Application with a default set of configurations.
func Default() *Application {
	return New(WithCORS(DefaultCORSConfig))
}

// Application is a web application.
type Application struct {
	decoder      Decoder
	port         int
	logger       Logger
	routeOptions []RouteOption

	inShutdown    int32
	readyCh       chan struct{}
	router        *httprouter.Router
	srvShutdownCh chan struct{}
	srv           *http.Server
	wg            sync.WaitGroup
}

// Run starts an HTTP server.
func (app *Application) Run() error {
	if app.shuttingDown() {
		return http.ErrServerClosed
	}

	defer app.wg.Wait()

	app.setupGracefulShutdown()

	err := app.listenAndServe()
	if err != nil && err != http.ErrServerClosed {
		app.logger.Println("Error:", err)
		app.shutdown()
		return err
	}

	return nil
}

// GET registers a new GET route for a path with handler.
func (app *Application) GET(path string, h Handler, opts ...RouteOption) {
	r := app.newRoute(h, opts)
	app.router.GET(path, r.buildHandle())
}

// POST registers a new POST route for a path with handler.
func (app *Application) POST(path string, h Handler, opts ...RouteOption) {
	r := app.newRoute(h, opts)
	app.router.POST(path, r.buildHandle())
}

// PUT registers a new PUT route for a path with handler.
func (app *Application) PUT(path string, h Handler, opts ...RouteOption) {
	r := app.newRoute(h, opts)
	app.router.PUT(path, r.buildHandle())
}

// PATCH registers a new PATCH route for a path with handler.
func (app *Application) PATCH(path string, h Handler, opts ...RouteOption) {
	r := app.newRoute(h, opts)
	app.router.PATCH(path, r.buildHandle())

}

// DELETE registers a new DELETE route for a path with handler.
func (app *Application) DELETE(path string, h Handler, opts ...RouteOption) {
	r := app.newRoute(h, opts)
	app.router.DELETE(path, r.buildHandle())
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

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (app *Application) listenAndServe() error {
	defer close(app.srvShutdownCh)

	app.srv = &http.Server{
		Handler: app.router,
		Addr:    fmt.Sprint(":", app.port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	addr := fmt.Sprint(":", app.port)
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	close(app.readyCh)
	app.logger.Println("Serving at port", app.port)
	return app.srv.Serve(ln)
}

// setupGracefulShutdown starts a goroutine for interrupt signal and proceed with graceful shutdown.
func (app *Application) setupGracefulShutdown() {
	app.execute(func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		select {
		case <-sigint:
			app.shutdown()
		case <-app.srvShutdownCh:
		}
	})
}

func (app *Application) shuttingDown() bool {
	return atomic.LoadInt32(&app.inShutdown) != 0
}

func (app *Application) shutdown() {
	atomic.StoreInt32(&app.inShutdown, 1)
	app.execute(func() {
		// We received an interrupt signal, shut down.
		if err := app.srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
			// Error from closing listeners, or context timeout:
			app.logger.Println("HTTP server Shutdown: %v", err)
		}
	})
}

func (app *Application) applyOpts(opts []Option) {
	for _, o := range opts {
		if o != nil {
			o.Apply(app)
		}
	}
}

// newRoute creates a new route give Handler and a list of RouteOption.
func (app *Application) newRoute(h Handler, opts []RouteOption) *route {
	r := &route{
		decoder: app.decoder,
		handler: h,
		logger:  app.logger,
	}

	r.applyOpts(app.routeOptions)
	r.applyOpts(opts)

	return r
}
