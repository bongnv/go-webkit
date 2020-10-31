package gwf

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

			bw, ok := rw.(*bufRespWriter)
			if !ok {
				next(rw, req, params)
				return
			}

			grw := newGzipRespWriter(cfg.Level, bw)
			next(grw, req, params)
			_ = grw.Close()
		}
	}
}

func newGzipRespWriter(level int, bw *bufRespWriter) *gzipResponseWriter {
	w, err := gzip.NewWriterLevel(bw, level)
	if err != nil {
		w = gzip.NewWriter(bw)
	}

	return &gzipResponseWriter{
		writer:        w,
		bufRespWriter: bw,
	}
}

type gzipResponseWriter struct {
	*bufRespWriter
	writer *gzip.Writer
}

func (gw gzipResponseWriter) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	return gw.writer.Write(b)
}

func (gw gzipResponseWriter) Close() error {
	header := gw.bufRespWriter.Header()
	header.Del(HeaderContentLength)
	if gw.statusCode == http.StatusNoContent {
		gw.writer.Reset(ioutil.Discard)
	} else {
		header.Set(HeaderContentEncoding, gzipScheme)
	}
	return gw.writer.Close()
}
