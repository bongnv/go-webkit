package webkit

import (
	"context"
	"net/http"
)

type contextKey int

const (
	_ contextKey = iota
	ctxKeyHTTPResponseWriter
)

// ResponseHeaderFromCtx returns Header for HTTP response which will be sent.
// It returns nil if it doesn't exist.
func ResponseHeaderFromCtx(ctx context.Context) http.Header {
	w, ok := ctx.Value(ctxKeyHTTPResponseWriter).(http.ResponseWriter)
	if ok {
		return w.Header()
	}

	return nil
}
