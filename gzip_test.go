package nanny

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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
		rw.WriteHeader(http.StatusOK)
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

	t.Run("using-gzip", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(HeaderAcceptEncoding, gzipScheme)
		rr := httptest.NewRecorder()
		h(rr, req, nil)
		require.Equal(t, http.StatusOK, rr.Code)
		require.Equal(t, "OK", decodeGzip(rr.Body))
		require.Equal(t, "gzip", rr.Header().Get(HeaderContentEncoding))
	})
}

func Test_WithGzip_fallback(t *testing.T) {
	var b bytes.Buffer
	app := &Application{
		logger: log.New(&b, "", log.LstdFlags),
	}

	h := func(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {}
	h = gzipTransformer(GzipConfig{
		Level: -200,
	})(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.WithValue(req.Context(), ctxKeyApp, app)
	req = req.WithContext(ctx)
	req.Header.Add(HeaderAcceptEncoding, gzipScheme)
	rr := httptest.NewRecorder()
	h(rr, req, nil)
	require.Equal(t, http.StatusOK, rr.Code)
	require.True(t, strings.Contains(b.String(), "Fallback"))
}

func Test_WithGzip_write_default_status(t *testing.T) {
	h := func(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {}
	h = gzipTransformer(DefaultGzipConfig)(h)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(HeaderAcceptEncoding, gzipScheme)
	rr := httptest.NewRecorder()
	h(rr, req, nil)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "", decodeGzip(rr.Body))
	require.Equal(t, "gzip", rr.Header().Get(HeaderContentEncoding))
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

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(HeaderAcceptEncoding, gzipScheme)
	rr := httptest.NewRecorder()

	h(rr, req, nil)

	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Equal(t, "", rr.Body.String())
	require.Equal(t, "", rr.Header().Get(HeaderContentEncoding))
}
