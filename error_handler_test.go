package nanny

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithErrorHandler(t *testing.T) {
	opt := WithErrorHandler(defaultErrorHandler())
	r := &route{}
	opt(r)
	require.NotNil(t, r.errorHandler)
}

func Test_defaultErrorHandler_CustomHTTPResponse(t *testing.T) {
	rr := httptest.NewRecorder()
	err := &HTTPError{
		Code:    http.StatusNotFound,
		Message: "Resource not found",
	}
	defaultErrorHandler()(rr, err)
	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Equal(t, "{\"message\":\"Resource not found\"}\n", rr.Body.String())
}

func Test_defaultErrorHandler_error(t *testing.T) {
	rr := httptest.NewRecorder()
	err := errors.New("resource not found")
	defaultErrorHandler()(rr, err)
	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "resource not found", rr.Body.String())
}
