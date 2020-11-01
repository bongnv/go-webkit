package gwf

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bongnv/inject"
	"github.com/julienschmidt/httprouter"
)

// Handler defines a function to serve HTTP requests.
type Handler func(ctx context.Context, req Request) (interface{}, error)

// New creates a new application.
func New(opts ...Option) *Application {
	logger := defaultLogger()
	app := &Application{
		container: inject.New(),
		router:    httprouter.New(),
		port:      8080,
		routeOptions: []RouteOption{
			WithDecoder(newDecoder()),
			WithEncoder(newEncoder()),
		},
		readyCh:       make(chan struct{}),
		srvShutdownCh: make(chan struct{}),
		logger:        logger,
	}

	app.applyOpts(opts)
	app.root = app.Group("/")
	return app
}

// Default returns an Application with a default set of configurations.
func Default() *Application {
	return New(
		WithLogger(defaultLogger()),
		WithRecovery(),
		WithCORS(DefaultCORSConfig),
		WithGzip(DefaultGzipConfig),
	)
}

// Application is a web application.
type Application struct {
	port         int
	logger       Logger
	routeOptions []RouteOption

	container     *inject.Container
	inShutdown    int32
	readyCh       chan struct{}
	root          *Group
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

// Component finds and returns a component via name.
// It returns an error if the requested component couldn't be found.
func (app *Application) Component(name string) (interface{}, error) {
	return app.container.Get(name)
}

// Register registers a new component to the application.
func (app *Application) Register(name string, component interface{}) error {
	return app.container.Register(name, component)
}

// MustRegister registers a new component to the application. It panics if there is any error.
func (app *Application) MustRegister(name string, component interface{}) {
	app.container.MustRegister(name, component)
}

// GET registers a new GET route for a path with handler.
func (app *Application) GET(path string, h Handler, opts ...RouteOption) {
	app.root.GET(path, h, opts...)
}

// POST registers a new POST route for a path with handler.
func (app *Application) POST(path string, h Handler, opts ...RouteOption) {
	app.root.POST(path, h, opts...)
}

// PUT registers a new PUT route for a path with handler.
func (app *Application) PUT(path string, h Handler, opts ...RouteOption) {
	app.root.PUT(path, h, opts...)
}

// PATCH registers a new PATCH route for a path with handler.
func (app *Application) PATCH(path string, h Handler, opts ...RouteOption) {
	app.root.PATCH(path, h, opts...)

}

// DELETE registers a new DELETE route for a path with handler.
func (app *Application) DELETE(path string, h Handler, opts ...RouteOption) {
	app.root.DELETE(path, h, opts...)
}

// Group creates a group of sub-routes
func (app *Application) Group(prefix string, opts ...RouteOption) *Group {
	if len(prefix) == 0 || prefix[0] != '/' {
		panic("path must begin with '/' in path '" + prefix + "'")
	}

	// Strip trailing / (if present) as all added sub paths must start with a /
	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	return &Group{
		prefix:       prefix,
		routeOptions: append(app.routeOptions, opts...),
		router:       app.router,
		logger:       app.logger,
	}
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
