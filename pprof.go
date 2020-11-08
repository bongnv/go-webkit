package nanny

import (
	"net/http"
	"net/http/pprof"
)

// WithPProf enables the pprof server.
func WithPProf(addr string) OptionFn {
	return func(app *Application) {
		mux := &http.ServeMux{}

		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		app.pprofSrv = &http.Server{
			Addr:    addr,
			Handler: mux,
		}
	}
}
