package webkit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WithDecoder(t *testing.T) {
	opt := WithDecoder(newDecoder())
	r := &route{}
	opt.ApplyRoute(r)
	require.NotNil(t, r.decoder)
}

func Test_decoderImpl(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/mock-request?age=10", strings.NewReader(`{"name": "mock-name"}`))
	req.Header.Add("Content-Type", "application/json")
	require.NoError(t, req.ParseForm())

	reqObj := &mockRequest{}
	d := newDecoder()
	err := d.Decode(reqObj, req)
	require.NoError(t, err)
	require.Equal(t, "mock-name", reqObj.Name)
	require.Equal(t, 10, reqObj.Age)
}
