package nanny

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/bongnv/inject"
	"github.com/julienschmidt/httprouter"
)

// Handler defines a function to serve HTTP requests.
type Handler func(ctx context.Context, req Request) (interface{}, error)

// DefaultApp is a plugin to provide a set of common options for an application.
var DefaultApp Plugin = []Option{
	WithLogger(defaultLogger()),
	WithRecovery(),
	WithCORS(DefaultCORSConfig),
	WithGzip(DefaultGzipConfig),
	WithTimeout(1 * time.Second),
	WithPProf(":8081"),
}

// New creates a new application.
func New(opts ...Option) *Application {
	app := &Application{
		addr:           ":8080",
		container:      inject.New(),
		readyCh:        make(chan struct{}),
		shutdownSignal: make(chan struct{}),
		logger:         defaultLogger(),
	}

	app.applyOpts([]Option{
		injectTimeoutMiddleware(),
		contextInjector(),
		WithDecoder(newDecoder()),
		WithEncoder(defaultEncoder{}),
	})

	app.applyOpts(opts)

	app.RouteGroup = &RouteGroup{
		routeOptions: app.routeOptions,
		app:          app,
	}

	return app
}

// Default returns an Application with a default set of configurations.
func Default(opts ...Option) *Application {
	opts = append(DefaultApp, opts...)
	return New(opts...)
}

// Application is a web application.
type Application struct {
	*RouteGroup

	addr         string
	logger       Logger
	routeOptions []RouteOption

	container *inject.Container
	readyCh   chan struct{}
	routes    []*route
	srv       *http.Server
	pprofSrv  *http.Server
	wg        sync.WaitGroup

	shutdownOnce   sync.Once
	shutdownSignal chan struct{}
}

// Run starts an HTTP server.
func (app *Application) Run() error {
	app.startHTTPServer()
	app.startPProfServer()
	app.setupGracefulShutdown()

	app.wg.Wait()
	return nil
}

// Component finds and returns a component via name.
// It returns an error if the requested component couldn't be found.
func (app *Application) Component(name string) (interface{}, error) {
	return app.container.Get(name)
}

// MustComponent finds and returns a component via name.
// It panics if there is any error.
func (app *Application) MustComponent(name string) interface{} {
	return app.container.MustGet(name)
}

// Register registers a new component to the application.
func (app *Application) Register(name string, component interface{}) error {
	return app.container.Register(name, component)
}

// MustRegister registers a new component to the application. It panics if there is any error.
func (app *Application) MustRegister(name string, component interface{}) {
	app.container.MustRegister(name, component)
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

func (app *Application) buildHTTPHandler() http.Handler {
	router := httprouter.New()

	for _, r := range app.routes {
		router.Handle(r.method, r.path, r.buildHandle())
	}

	return router
}

// ListenAndServe listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (app *Application) startHTTPServer() {
	app.execute(func() {
		app.srv = &http.Server{
			Addr:    app.addr,
			Handler: app.buildHTTPHandler(),
		}

		ln, err := net.Listen("tcp", app.addr)
		if err != nil {
			app.logger.Println("Error when listening on", app.addr)
			return
		}

		close(app.readyCh)
		app.logger.Println("Serving at addr", app.addr)
		if err := app.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			app.logger.Println("Error when starting HTTP service", err)
		}

		app.shutdown()
	})
}

func (app *Application) startPProfServer() {
	if app.pprofSrv == nil {
		return
	}

	app.execute(func() {
		if err := app.pprofSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Println("Error while starting pprof service", err)
		}
		app.shutdown()
	})
}

// setupGracefulShutdown starts a goroutine for interrupt signal and proceed with graceful shutdown.
func (app *Application) setupGracefulShutdown() {
	app.execute(func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		select {
		case <-sigint:
		case <-app.shutdownSignal:
		}

		app.execute(func() {
			// We received an interrupt signal, shut down.
			if err := app.srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
				// Error from closing listeners, or context timeout:
				app.logger.Println("HTTP server Shutdown: %v", err)
			}
		})

		if app.pprofSrv != nil {
			app.execute(func() {
				_ = app.pprofSrv.Shutdown(context.Background())
			})
		}
	})
}

func (app *Application) shutdown() {
	app.shutdownOnce.Do(func() {
		close(app.shutdownSignal)
	})
}

func (app *Application) applyOpts(opts []Option) {
	for _, o := range opts {
		if o != nil {
			o.Apply(app)
		}
	}
}
