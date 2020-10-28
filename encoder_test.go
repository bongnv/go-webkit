package webkit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockResponse struct {
	Data string `json:"data"`
}

func Test_defaultEncoder(t *testing.T) {
	e := &defaultEncoder{}
	rr := httptest.NewRecorder()
	err := e.Encode(rr, &mockResponse{Data: "mock-data"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "{\"data\":\"mock-data\"}\n", rr.Body.String())
}
