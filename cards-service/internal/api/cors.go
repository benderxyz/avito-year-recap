package api

import (
	"net/http"
	"strings"
)

const corsPreflightMaxAge = "600"

// WithCORS answers cross-origin requests from the configured origins and
// short-circuits preflight requests. Origins are matched exactly; a request
// from an unlisted origin is served without CORS headers, so the browser
// blocks it.
func WithCORS(next http.Handler, allowedOrigins []string) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Vary on Origin regardless of the outcome: the response body is the
		// same for every origin, but the CORS headers are not, and caches
		// must not reuse an allowed origin's response for a denied one.
		w.Header().Add("Vary", "Origin")

		if _, ok := allowed[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Add("Vary", "Access-Control-Request-Method")
			w.Header().Add("Vary", "Access-Control-Request-Headers")

			if _, ok := allowed[origin]; ok {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(corsAllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsAllowedHeaders, ", "))
				w.Header().Set("Access-Control-Max-Age", corsPreflightMaxAge)
			}

			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

var (
	corsAllowedMethods = []string{http.MethodGet, http.MethodOptions}
	corsAllowedHeaders = []string{"Content-Type"}
)
