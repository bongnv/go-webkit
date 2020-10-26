package webkit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var mockRouteOptionFn RouteOptionFn = func(r *route) {
	r.middlewares = append(r.middlewares, mockMiddleware)
}

func Test_RouteOptionFn_ApplyRoute(t *testing.T) {
	r := &route{}
	mockRouteOptionFn.ApplyRoute(r)
	require.Len(t, r.middlewares, 1)
}

func Test_RouteOptionFn_Apply(t *testing.T) {
	app := &Application{}
	mockRouteOptionFn.Apply(app)
	require.Len(t, app.routeOptions, 1)
}

func Test_buildHandle_responseError(t *testing.T) {
	r := &route{
		handler: func(_ context.Context, _ Request) error {
			return errors.New("remote error")
		},
	}
	rr := httptest.NewRecorder()
	handle := r.buildHandle()
	handle(rr, &http.Request{}, nil)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.True(t, strings.Contains(rr.Body.String(), "remote error"))
}

func Test_route_applyOpts(t *testing.T) {
	r := &route{}
	r.applyOpts([]RouteOption{mockMiddleware})
	require.Len(t, r.middlewares, 1)
}
