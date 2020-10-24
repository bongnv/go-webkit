package webkit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

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
		return nil
	})
	req, _ := http.NewRequest("GET", "/mock-endpoint", nil)
	resp := executeRequest(app, req)
	require.Equal(t, http.StatusOK, resp.Code)
}
