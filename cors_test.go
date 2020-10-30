package gwf

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func mockHandler(ctx context.Context, req Request) (interface{}, error) {
	return nil, nil
}

func mockContext(w http.ResponseWriter) context.Context {
	return context.WithValue(context.Background(), ctxKeyHTTPResponseWriter, w)
}

func TestCORS_wildcard_origin(t *testing.T) {
	rr := httptest.NewRecorder()
	ctx := mockContext(rr)
	req := &requestImpl{
		httpReq: httptest.NewRequest(http.MethodGet, "/", nil),
	}
	h := WithCORS(DefaultCORSConfig)(mockHandler)
	_, err := h(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "*", rr.Header().Get(HeaderAccessControlAllowOrigin))
}

func TestCORS_allow_origin(t *testing.T) {
	rr := httptest.NewRecorder()
	ctx := mockContext(rr)
	req := &requestImpl{
		httpReq: httptest.NewRequest(http.MethodGet, "/", nil),
	}
	h := WithCORS(CORSConfig{
		AllowOrigins:     []string{"localhost"},
		AllowCredentials: true,
	})(mockHandler)
	req.httpReq.Header.Set(HeaderOrigin, "localhost")
	_, err := h(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "localhost", rr.Header().Get(HeaderAccessControlAllowOrigin))
	require.Equal(t, "true", rr.Header().Get(HeaderAccessControlAllowCredentials))
}

func TestCORS_preflight_request(t *testing.T) {
	rr := httptest.NewRecorder()
	ctx := mockContext(rr)
	req := &requestImpl{
		httpReq: httptest.NewRequest(http.MethodOptions, "/", nil),
	}

	req.httpReq.Header.Set(HeaderOrigin, "localhost")
	req.httpReq.Header.Set(HeaderAccessControlRequestHeaders, "Content-Type")
	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"localhost"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet},
		MaxAge:           3600,
	})
	h := cors(mockHandler)
	_, err := h(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "localhost", rr.Header().Get(HeaderAccessControlAllowOrigin))
	require.NotEmpty(t, rr.Header().Get(HeaderAccessControlAllowMethods))
	require.Equal(t, "true", rr.Header().Get(HeaderAccessControlAllowCredentials))
	require.Equal(t, "3600", rr.Header().Get(HeaderAccessControlMaxAge))
	require.Equal(t, "Content-Type", rr.Header().Get(HeaderAccessControlAllowHeaders))
}

func TestCORS_preflight_with_wildcard(t *testing.T) {
	// Preflight request with `AllowOrigins` *
	rr := httptest.NewRecorder()
	ctx := mockContext(rr)
	req := &requestImpl{
		httpReq: httptest.NewRequest(http.MethodOptions, "/", nil),
	}
	req.httpReq.Header.Set(HeaderOrigin, "localhost")

	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet},
		AllowHeaders:     []string{HeaderContentType},
		MaxAge:           3600,
	})
	h := cors(mockHandler)
	_, err := h(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "localhost", rr.Header().Get(HeaderAccessControlAllowOrigin))
	require.NotEmpty(t, rr.Header().Get(HeaderAccessControlAllowMethods))
	require.Equal(t, "true", rr.Header().Get(HeaderAccessControlAllowCredentials))
	require.Equal(t, "3600", rr.Header().Get(HeaderAccessControlMaxAge))
	require.Equal(t, "Content-Type", rr.Header().Get(HeaderAccessControlAllowHeaders))
}

func TestCORS_preflight_with_subdomain(t *testing.T) {
	// Preflight request with `AllowOrigins` which allow all subdomains with *
	rr := httptest.NewRecorder()
	ctx := mockContext(rr)
	req := &requestImpl{
		httpReq: httptest.NewRequest(http.MethodOptions, "/", nil),
	}
	req.httpReq.Header.Set(HeaderOrigin, "https://a.example.com")
	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"https://*.example.com"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h := cors(mockHandler)
	_, err := h(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "https://a.example.com", rr.Header().Get(HeaderAccessControlAllowOrigin))

	req.httpReq.Header.Set(HeaderOrigin, "https://b.example.com")
	_, err = h(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "https://b.example.com", rr.Header().Get(HeaderAccessControlAllowOrigin))
}

func Test_CORS_allowOriginScheme(t *testing.T) {
	rr := httptest.NewRecorder()
	ctx := mockContext(rr)
	req := &requestImpl{
		httpReq: httptest.NewRequest(http.MethodOptions, "/", nil),
	}
	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"https://example.com"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h := cors(mockHandler)

	req.httpReq.Header.Set(HeaderOrigin, "https://example.com")
	_, err := h(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "https://example.com", rr.Header().Get(HeaderAccessControlAllowOrigin))

	req.httpReq.Header.Set(HeaderOrigin, "http://example.com")
	_, err = h(ctx, req)
	require.NoError(t, err)
	require.Empty(t, rr.Header().Get(HeaderAccessControlAllowOrigin))
}
