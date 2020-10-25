package webkit

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_buildHandlerFunc_responseError(t *testing.T) {
	rr := httptest.NewRecorder()
	handle := buildHandlerFunc(func(_ context.Context, _ Request) error {
		return errors.New("remote error")
	})
	handle(rr, &http.Request{}, nil)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.True(t, strings.Contains(rr.Body.String(), "remote error"))
}
