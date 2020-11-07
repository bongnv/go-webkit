package nanny

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var mockMiddleware Middleware = func(h Handler) Handler {
	return h
}

func Test_ApplyRoute(t *testing.T) {
	r := &route{}
	mockMiddleware.ApplyRoute(r)
	require.Len(t, r.middlewares, 1)
}

func Test_Apply(t *testing.T) {
	a := &Application{}
	mockMiddleware.Apply(a)
	require.Len(t, a.routeOptions, 1)
}
