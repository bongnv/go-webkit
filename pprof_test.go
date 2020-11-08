package nanny

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithPProf(t *testing.T) {
	app := &Application{}
	WithPProf(":8081")(app)
	require.NotNil(t, app.pprofSrv)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/", nil)
	app.pprofSrv.Handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
