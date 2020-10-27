package webkit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

type mockResponse struct {
	Data string `json:"data"`
}

type mockRequest struct {
	Name string
	Age  int
}

func Test_Decode(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/mock-request?age=10", nil)
	r := &requestImpl{
		decoder: newDecoder(),
		httpReq: req,
		params: []httprouter.Param{{
			Key:   "name",
			Value: "mock-name",
		}},
	}
	reqObj := &mockRequest{}
	err := r.Decode(reqObj)
	require.NoError(t, err)
	require.Equal(t, "mock-name", reqObj.Name)
	require.Equal(t, 10, reqObj.Age)
}

func Test_Response(t *testing.T) {
	rr := httptest.NewRecorder()
	r := &requestImpl{
		httpWriter: rr,
	}

	err := r.Respond(&mockResponse{Data: "mock-data"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "{\"data\":\"mock-data\"}\n", rr.Body.String())
}
