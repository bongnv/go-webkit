package gwf

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HTTPError(t *testing.T) {
	err := &HTTPError{
		Code:    http.StatusNotFound,
		Message: "Some error message",
	}

	require.EqualError(t, err, "Some error message")

	code, body := err.HTTPResponse()
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, `{"message":"Some error message"}`, string(body))
}
