package admin

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

const bearerPrefix = "Bearer "

func RequireToken(expected string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !tokenMatches(expected, r.Header.Get("Authorization")) {
			writeError(w, http.StatusUnauthorized, AdminError{Error: "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func tokenMatches(expected, header string) bool {
	if expected == "" {
		return false
	}
	if len(header) <= len(bearerPrefix) || !strings.EqualFold(header[:len(bearerPrefix)], bearerPrefix) {
		return false
	}

	provided := strings.TrimSpace(header[len(bearerPrefix):])

	return subtle.ConstantTimeCompare([]byte(provided), []byte(expected)) == 1
}
