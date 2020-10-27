package webkit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func mockHandler(ctx context.Context, req Request) error {
	return nil
}

func TestCORS_wildcard_origin(t *testing.T) {
	rr := httptest.NewRecorder()
	req := &requestImpl{
		httpReq:    httptest.NewRequest(http.MethodGet, "/", nil),
		httpWriter: rr,
	}
	h := WithCORS(DefaultCORSConfig)(mockHandler)
	require.NoError(t, h(context.Background(), req))
	require.Equal(t, "*", rr.Header().Get(HeaderAccessControlAllowOrigin))
}

func TestCORS_allow_origin(t *testing.T) {
	rr := httptest.NewRecorder()
	req := &requestImpl{
		httpReq:    httptest.NewRequest(http.MethodGet, "/", nil),
		httpWriter: rr,
	}
	h := WithCORS(CORSConfig{
		AllowOrigins: []string{"localhost"},
	})(mockHandler)
	req.httpReq.Header.Set(HeaderOrigin, "localhost")
	require.NoError(t, h(context.Background(), req))
	require.Equal(t, "localhost", rr.Header().Get(HeaderAccessControlAllowOrigin))
}

func TestCORS_preflight_request(t *testing.T) {
	rr := httptest.NewRecorder()
	req := &requestImpl{
		httpReq:    httptest.NewRequest(http.MethodOptions, "/", nil),
		httpWriter: rr,
	}

	req.httpReq.Header.Set(HeaderOrigin, "localhost")
	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"localhost"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet},
		MaxAge:           3600,
	})
	h := cors(mockHandler)
	require.NoError(t, h(context.Background(), req))
	require.Equal(t, "localhost", rr.Header().Get(HeaderAccessControlAllowOrigin))
	require.NotEmpty(t, rr.Header().Get(HeaderAccessControlAllowMethods))
	require.Equal(t, "true", rr.Header().Get(HeaderAccessControlAllowCredentials))
	require.Equal(t, "3600", rr.Header().Get(HeaderAccessControlMaxAge))
}

func TestCORS_preflight_with_wildcard(t *testing.T) {
	// Preflight request with `AllowOrigins` *
	rr := httptest.NewRecorder()
	req := &requestImpl{
		httpReq:    httptest.NewRequest(http.MethodOptions, "/", nil),
		httpWriter: rr,
	}
	req.httpReq.Header.Set(HeaderOrigin, "localhost")

	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet},
		MaxAge:           3600,
	})
	h := cors(mockHandler)
	require.NoError(t, h(context.Background(), req))
	require.Equal(t, "localhost", rr.Header().Get(HeaderAccessControlAllowOrigin))
	require.NotEmpty(t, rr.Header().Get(HeaderAccessControlAllowMethods))
	require.Equal(t, "true", rr.Header().Get(HeaderAccessControlAllowCredentials))
	require.Equal(t, "3600", rr.Header().Get(HeaderAccessControlMaxAge))
}

func TestCORS_preflight_with_subdomain(t *testing.T) {
	// Preflight request with `AllowOrigins` which allow all subdomains with *
	rr := httptest.NewRecorder()
	req := &requestImpl{
		httpReq:    httptest.NewRequest(http.MethodOptions, "/", nil),
		httpWriter: rr,
	}
	req.httpReq.Header.Set(HeaderOrigin, "https://a.example.com")
	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"https://*.example.com"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h := cors(mockHandler)
	require.NoError(t, h(context.Background(), req))
	require.Equal(t, "https://a.example.com", rr.Header().Get(HeaderAccessControlAllowOrigin))

	req.httpReq.Header.Set(HeaderOrigin, "https://b.example.com")
	require.NoError(t, h(context.Background(), req))
	require.Equal(t, "https://b.example.com", rr.Header().Get(HeaderAccessControlAllowOrigin))
}

func Test_CORS_allowOriginScheme(t *testing.T) {
	rr := httptest.NewRecorder()
	req := &requestImpl{
		httpReq:    httptest.NewRequest(http.MethodOptions, "/", nil),
		httpWriter: rr,
	}
	cors := WithCORS(CORSConfig{
		AllowOrigins:     []string{"https://example.com"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h := cors(mockHandler)

	req.httpReq.Header.Set(HeaderOrigin, "https://example.com")
	require.NoError(t, h(context.Background(), req))
	require.Equal(t, "https://example.com", rr.Header().Get(HeaderAccessControlAllowOrigin))

	req.httpReq.Header.Set(HeaderOrigin, "http://example.com")
	require.NoError(t, h(context.Background(), req))
	require.Empty(t, rr.Header().Get(HeaderAccessControlAllowOrigin))
}
