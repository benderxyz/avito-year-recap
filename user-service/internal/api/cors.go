package api

import (
	"net/http"
	"strings"
)

const corsPreflightMaxAge = "600"

func WithCORS(next http.Handler, allowedOrigins []string) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		allowed[origin] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

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
