package nanny

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Recovery(t *testing.T) {
	var b bytes.Buffer
	r := &route{
		handler: func(_ context.Context, req Request) (interface{}, error) {
			panic("random panic")
		},
		logger:       log.New(&b, "", log.LstdFlags),
		errorHandler: defaultErrorHandler(),
		middlewares:  nil,
	}

	WithRecovery().ApplyRoute(r)
	handle := r.buildHandle()
	rr := httptest.NewRecorder()

	require.NotPanics(t, func() {
		handle(rr, &http.Request{}, nil)
	})

	require.Equal(t, http.StatusServiceUnavailable, rr.Code)
	require.Equal(t, "{\"message\":\"Service Unavailable\"}\n", rr.Body.String())
	require.True(t, strings.Contains(b.String(), "[PANIC RECOVER]"))
	require.True(t, strings.Contains(b.String(), "random panic"))
}
