package nanny

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"
)

func Test_WithGzip(t *testing.T) {
	opt := WithGzip(DefaultGzipConfig)
	r := &route{}
	opt(r)
	require.Len(t, r.transformers, 1)

	h := func(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		_, _ = rw.Write([]byte("OK"))
	}
	h = r.transformers[0](h)

	t.Run("ignore-by-header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		h(rr, req, nil)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "OK", rr.Body.String())
		require.Equal(t, "", rr.Header().Get(HeaderContentEncoding))
	})

	t.Run("ignore-by-type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(HeaderAcceptEncoding, gzipScheme)
		rr := httptest.NewRecorder()
		h(rr, req, nil)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "OK", rr.Body.String())
		require.Equal(t, "", rr.Header().Get(HeaderContentEncoding))
	})

	t.Run("using-gzip", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(HeaderAcceptEncoding, gzipScheme)
		rr := httptest.NewRecorder()
		hWithBuf := brwTransformer(h)
		hWithBuf(rr, req, nil)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "OK", decodeGzip(rr.Body))
		require.Equal(t, "gzip", rr.Header().Get(HeaderContentEncoding))
	})
}

func decodeGzip(body io.Reader) string {
	gr, _ := gzip.NewReader(body)
	s, _ := ioutil.ReadAll(gr)
	return string(s)
}

func Test_WithGzip_no_content(t *testing.T) {
	h := func(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		rw.WriteHeader(http.StatusNoContent)
	}

	h = gzipTransformer(DefaultGzipConfig)(h)
	h = brwTransformer(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(HeaderAcceptEncoding, gzipScheme)
	rr := httptest.NewRecorder()

	h(rr, req, nil)

	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Equal(t, "", rr.Body.String())
	require.Equal(t, "", rr.Header().Get(HeaderContentEncoding))
}
