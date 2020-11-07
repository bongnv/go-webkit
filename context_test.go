package nanny

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResponseHeaderFromCtx(t *testing.T) {
	rr := httptest.NewRecorder()
	ctx := context.WithValue(context.Background(), ctxKeyHTTPResponseWriter, rr)
	require.NotNil(t, ResponseHeaderFromCtx(ctx))
}

func Test_ResponseHeaderFromCtx_nil(t *testing.T) {
	require.Nil(t, ResponseHeaderFromCtx(context.Background()))
}
