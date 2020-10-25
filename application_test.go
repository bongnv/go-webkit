package webkit

import (
	"context"
	"net/http"
	"net/http/httptest"
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
	app.GET("/mock-endpoint", func(ctx context.Context, req Request) error {
		return req.Response("OK")
	})
	req, _ := http.NewRequest("GET", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
	require.Equal(t, "\"OK\"\n", resp.Body.String())
}

func Test_PUT(t *testing.T) {
	app := New()
	app.PUT("/mock-endpoint", func(ctx context.Context, req Request) error {
		return nil
	})
	req, _ := http.NewRequest("PUT", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
}

func Test_POST(t *testing.T) {
	app := New()
	app.POST("/mock-endpoint", func(ctx context.Context, req Request) error {
		return nil
	})
	req, _ := http.NewRequest("POST", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
}

func Test_DELETE(t *testing.T) {
	app := New()
	app.DELETE("/mock-endpoint", func(ctx context.Context, req Request) error {
		return nil
	})
	req, _ := http.NewRequest("DELETE", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
}

func Test_PATCH(t *testing.T) {
	app := New()
	app.PATCH("/mock-endpoint", func(ctx context.Context, req Request) error {
		return nil
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
