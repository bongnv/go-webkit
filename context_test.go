package nanny

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

func Test_ResponseHeaderFromCtx(t *testing.T) {
	rr := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), ctxKeyHTTPResponseWriter, rr)
	require.NotNil(t, ResponseHeaderFromCtx(ctx))
}

func Test_ResponseHeaderFromCtx_nil(t *testing.T) {
	require.Nil(t, ResponseHeaderFromCtx(context.Background()))
}

func Test_contextInjector(t *testing.T) {
	app := &Application{}
	contextInjector().Apply(app)
	require.Len(t, app.routeOptions, 1)
	r := &route{}
	app.routeOptions[0].ApplyRoute(r)
	require.Len(t, r.transformers, 1)

	h := func(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
		a := req.Context().Value(ctxKeyApp)
		require.NotNil(t, a)
		require.IsType(t, &Application{}, a)

		rr := req.Context().Value(ctxKeyHTTPResponseWriter)
		require.NotNil(t, rr)

		r := req.Context().Value(ctxKeyRoute)
		require.NotNil(t, r)
	}
	h = r.transformers[0](h)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	h(rr, req, nil)
}

func Test_loggerFromCtx(t *testing.T) {
	app := &Application{
		logger: defaultLogger(),
	}

	ctx := context.WithValue(context.Background(), ctxKeyApp, app)
	logger := loggerFromCtx(ctx)
	require.NotNil(t, logger)
}

func Test_loggerFromCtx_nil(t *testing.T) {
	logger := loggerFromCtx(context.Background())
	require.Nil(t, logger)
}
