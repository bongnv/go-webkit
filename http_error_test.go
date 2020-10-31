package gwf

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HTTPError(t *testing.T) {
	err := &HTTPError{
		Code:    http.StatusNotFound,
		Message: "Some error message",
	}

	require.EqualError(t, err, "Some error message")
	rr := httptest.NewRecorder()
	err.WriteTo(rr)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Equal(t, "{\"message\":\"Some error message\"}\n", rr.Body.String())
}
