package webkit

import (
	"context"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig is the default configuration for the CORS middleware.
var DefaultCORSConfig = CORSConfig{
	AllowOrigins: []string{"*"},
	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
}

// WithCORS returns a middleware to support Cross-Origin Resource Sharing.
func WithCORS(cfg CORSConfig) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req Request) error {
			allowMethods := strings.Join(cfg.AllowMethods, ",")
			allowHeaders := strings.Join(cfg.AllowHeaders, ",")

			httpReq := req.HTTPRequest()
			origin := httpReq.Header.Get(HeaderOrigin)
			allowOrigin := getAllowOrigin(origin, cfg)

			// non-OPTIONS requests
			if httpReq.Method != http.MethodOptions {
				req.ResponseHeader().Add(HeaderVary, HeaderOrigin)
				req.ResponseHeader().Set(HeaderAccessControlAllowOrigin, allowOrigin)
				if cfg.AllowCredentials {
					req.ResponseHeader().Set(HeaderAccessControlAllowCredentials, "true")
				}
				return next(ctx, req)
			}

			// Preflight requests
			req.ResponseHeader().Add(HeaderVary, HeaderOrigin)
			req.ResponseHeader().Add(HeaderVary, HeaderAccessControlRequestMethod)
			req.ResponseHeader().Add(HeaderVary, HeaderAccessControlRequestHeaders)
			req.ResponseHeader().Set(HeaderAccessControlAllowOrigin, allowOrigin)
			req.ResponseHeader().Set(HeaderAccessControlAllowMethods, allowMethods)
			if cfg.AllowCredentials {
				req.ResponseHeader().Set(HeaderAccessControlAllowCredentials, "true")
			}
			if allowHeaders != "" {
				req.ResponseHeader().Set(HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				h := req.ResponseHeader().Get(HeaderAccessControlRequestHeaders)
				if h != "" {
					req.ResponseHeader().Set(HeaderAccessControlAllowHeaders, h)
				}
			}
			if cfg.MaxAge > 0 {
				req.ResponseHeader().Set(HeaderAccessControlMaxAge, strconv.Itoa(cfg.MaxAge))
			}

			return req.Respond(nil)
		}
	}
}

func matchScheme(domain, pattern string) bool {
	didx := strings.Index(domain, ":")
	pidx := strings.Index(pattern, ":")
	return didx != -1 && pidx != -1 && domain[:didx] == pattern[:pidx]
}

// matchSubdomain compares authority with wildcard
func matchSubdomain(domain, pattern string) bool {
	if !matchScheme(domain, pattern) {
		return false
	}
	didx := strings.Index(domain, "://")
	pidx := strings.Index(pattern, "://")
	if didx == -1 || pidx == -1 {
		return false
	}
	domAuth := domain[didx+3:]
	// to avoid long loop by invalid long domain
	if len(domAuth) > 253 {
		return false
	}
	patAuth := pattern[pidx+3:]

	domComp := strings.Split(domAuth, ".")
	patComp := strings.Split(patAuth, ".")
	for i := len(domComp)/2 - 1; i >= 0; i-- {
		opp := len(domComp) - 1 - i
		domComp[i], domComp[opp] = domComp[opp], domComp[i]
	}
	for i := len(patComp)/2 - 1; i >= 0; i-- {
		opp := len(patComp) - 1 - i
		patComp[i], patComp[opp] = patComp[opp], patComp[i]
	}

	for i, v := range domComp {
		if len(patComp) <= i {
			return false
		}
		p := patComp[i]
		if p == "*" {
			return true
		}
		if p != v {
			return false
		}
	}
	return false
}

func getAllowOrigin(origin string, cfg CORSConfig) string {
	// Check allowed origins
	for _, o := range cfg.AllowOrigins {
		if o == "*" && cfg.AllowCredentials {
			return origin
		}

		if o == "*" || o == origin {
			return o
		}

		if matchSubdomain(origin, o) {
			return origin
		}
	}

	var allowOriginPatterns []string
	for _, allowOrigin := range cfg.AllowOrigins {
		pattern := regexp.QuoteMeta(allowOrigin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		pattern = "^" + pattern + "$"
		allowOriginPatterns = append(allowOriginPatterns, pattern)
	}

	// Check allowed origin patterns
	for _, re := range allowOriginPatterns {
		didx := strings.Index(origin, "://")
		if didx == -1 {
			continue
		}
		domAuth := origin[didx+3:]
		// to avoid regex cost by invalid long domain
		if len(domAuth) > 253 {
			break
		}

		if match, _ := regexp.MatchString(re, origin); match {
			return origin
		}
	}

	return ""
}
