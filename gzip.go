package nanny

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

const (
	gzipScheme = "gzip"
)

// GzipConfig defines the config for Gzip middleware.
type GzipConfig struct {
	// Gzip compression level.
	// Optional. Default value -1.
	Level int
}

var (
	// DefaultGzipConfig is the default config for Gzip middleware.
	DefaultGzipConfig = GzipConfig{
		Level: gzip.DefaultCompression,
	}
)

// WithGzip returns a middleware which compresses HTTP response using gzip compression.
func WithGzip(cfg GzipConfig) RouteOptionFn {
	return func(r *route) {
		r.transformers = append(r.transformers, gzipTransformer(cfg))
	}
}

func gzipTransformer(cfg GzipConfig) handleTransformer {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
			if !strings.Contains(req.Header.Get(HeaderAcceptEncoding), gzipScheme) {
				next(rw, req, params)
				return
			}

			logger := loggerFromCtx(req.Context())
			w, err := gzip.NewWriterLevel(rw, cfg.Level)
			if err != nil {
				// TODO: There should be log
				logger.Println("Fallback to default compress due to", err)
				w = gzip.NewWriter(rw)
			}

			grw := &gzipResponseWriter{
				writer:         w,
				ResponseWriter: rw,
			}

			next(grw, req, params)
			if err := grw.writer.Close(); err != nil {
				logger.Println("Error while closing gzip writer", err)
			}
		}
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer     *gzip.Writer
	statusCode int
}

func (gw *gzipResponseWriter) WriteHeader(statusCode int) {
	gw.statusCode = statusCode

	if gw.isBodyAllowed() {
		header := gw.ResponseWriter.Header()
		header.Del(HeaderContentLength)
		header.Set(HeaderContentEncoding, gzipScheme)
	} else {
		gw.writer.Reset(ioutil.Discard)
	}

	gw.ResponseWriter.WriteHeader(statusCode)
}

func (gw *gzipResponseWriter) Write(b []byte) (int, error) {
	if gw.statusCode == 0 {
		gw.WriteHeader(http.StatusOK)
	}

	return gw.writer.Write(b)
}

func (gw gzipResponseWriter) isBodyAllowed() bool {
	return gw.statusCode != http.StatusNoContent
}
