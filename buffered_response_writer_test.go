package webkit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

func Test_brwTransformer_no_content(t *testing.T) {
	h := func(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		rw.WriteHeader(http.StatusNoContent)
		_, _ = rw.Write(nil)
	}
	h = brwTransformer(h)
	rr := httptest.NewRecorder()
	h(rr, nil, nil)
	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Empty(t, rr.Body.String())
}
