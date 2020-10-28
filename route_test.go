package webkit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
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

func Test_buildHandle(t *testing.T) {
	r := &route{
		encoder: newEncoder(),
		handler: func(_ context.Context, req Request) (interface{}, error) {
			require.IsType(t, &requestImpl{}, req)
			require.Len(t, req.(*requestImpl).params, 1)
			require.Equal(t, "key", req.(*requestImpl).params[0].Key)
			return "OK", nil
		},
	}
	rr := httptest.NewRecorder()
	handle := r.buildHandle()
	handle(rr, &http.Request{}, []httprouter.Param{{
		Key:   "key",
		Value: "value",
	}})

	require.Equal(t, http.StatusOK, rr.Code)
}

func Test_buildHandle_responseError(t *testing.T) {
	r := &route{
		handler: func(_ context.Context, _ Request) (interface{}, error) {
			return nil, errors.New("remote error")
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
