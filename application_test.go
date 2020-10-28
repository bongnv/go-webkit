package webkit

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
	app.router.ServeHTTP(rr, req)

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
	testDone := make(chan struct{})

	require.NotPanics(t, func() {
		defer close(testDone)
		app := New()
		runFinished := make(chan struct{})

		go func() {
			require.NoError(t, app.Run())
			close(runFinished)
		}()

		<-app.readyCh
		app.shutdown()
		<-runFinished

		select {
		case _, ok := <-app.srvShutdownCh:
			require.False(t, ok, "srvShutdownCh should be closed")
		default:
			require.Fail(t, "srvShutdownCh should be closed")
		}

		require.EqualError(t, app.Run(), http.ErrServerClosed.Error())
	})

	select {
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Test times out")
	case <-testDone:
	}
}

func Test_applyOpts(t *testing.T) {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	opt := WithLogger(logger)
	app := &Application{}
	app.applyOpts([]Option{opt})
	require.Equal(t, logger, app.logger)
}

func Test_Default(t *testing.T) {
	app := Default()
	require.Len(t, app.routeOptions, 2)
	require.NotNil(t, app.decoder)
	require.NotNil(t, app.router)
	require.NotNil(t, app.logger)
}
