package nanny

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func executeRequest(app *Application, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.buildHTTPHandler().ServeHTTP(rr, req)

	return rr
}

func Test_GET(t *testing.T) {
	app := New()
	app.GET("/mock-endpoint", func(ctx context.Context, req Request) (interface{}, error) {
		return "OK", nil
	})
	req, _ := http.NewRequest("GET", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "\"OK\"\n", resp.Body.String())
}

func Test_PUT(t *testing.T) {
	app := New()
	app.PUT("/mock-endpoint", func(ctx context.Context, req Request) (interface{}, error) {
		return "OK", nil
	})
	req, _ := http.NewRequest("PUT", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
}

func Test_POST(t *testing.T) {
	app := New()
	app.POST("/mock-endpoint", func(ctx context.Context, req Request) (interface{}, error) {
		return "OK", nil
	})
	req, _ := http.NewRequest("POST", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
}

func Test_DELETE(t *testing.T) {
	app := New()
	app.DELETE("/mock-endpoint", func(ctx context.Context, req Request) (interface{}, error) {
		return nil, nil
	})
	req, _ := http.NewRequest("DELETE", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusNoContent, resp.Code)
}

func Test_PATCH(t *testing.T) {
	app := New()
	app.PATCH("/mock-endpoint", func(ctx context.Context, req Request) (interface{}, error) {
		return "OK", nil
	})
	req, _ := http.NewRequest("PATCH", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
}

func Test_graceful_shutdown(t *testing.T) {
	runFinished := make(chan struct{})

	app := New()

	go func() {
		require.NotPanics(t, func() {
			app.Run()
			close(runFinished)
		})
	}()

	select {
	case <-app.readyCh:
		app.shutdown()
	case <-runFinished:
		require.Fail(t, "App shouldn't stop running")
	}

	select {
	case <-time.After(200 * time.Millisecond):
		require.Fail(t, "Test times out")
	case <-runFinished:
	}
}

func Test_applyOpts(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	opt := WithLogger(logger)
	app := New()
	app.applyOpts([]Option{opt})
	require.Equal(t, logger, app.logger)
}

func Test_Default(t *testing.T) {
	app := Default()
	require.Len(t, app.routeOptions, 8)
	require.NotNil(t, app.logger)
	require.NotNil(t, app.pprofSrv)
}

func Test_Application_Group(t *testing.T) {
	app := New()
	require.Panics(t, func() {
		app.Group("")
	}, "must panic with empty path")

	require.Panics(t, func() {
		app.Group("invalid")
	}, "must panic with invalid path")

	require.NotPanics(t, func() {
		g := app.Group("/")
		require.Empty(t, g.prefix)
	})
}

func Test_Component(t *testing.T) {
	app := New()
	require.NoError(t, app.Register("logger", defaultLogger()))
	l, err := app.Component("logger")
	require.NoError(t, err)
	require.NotNil(t, l)
	_, ok := l.(Logger)
	require.True(t, ok)
	require.NotNil(t, app.MustComponent("logger"))
}
