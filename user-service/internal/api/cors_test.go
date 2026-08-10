package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newCORSHandler() http.Handler {
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return WithCORS(next, []string{"http://localhost:3000"})
}

func TestWithCORSShouldAllowConfiguredOrigin(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/users", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()

	newCORSHandler().ServeHTTP(recorder, req)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("expected allow-origin for configured origin, got %q", got)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected request to reach next handler, got %d", recorder.Code)
	}
}

func TestWithCORSShouldNotAllowUnknownOrigin(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/users", nil)
	req.Header.Set("Origin", "http://evil.example")
	recorder := httptest.NewRecorder()

	newCORSHandler().ServeHTTP(recorder, req)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allow-origin for unknown origin, got %q", got)
	}
}

func TestWithCORSShouldAnswerPreflight(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/users", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	recorder := httptest.NewRecorder()

	newCORSHandler().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for preflight, got %d", recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Fatal("expected allow-methods on preflight response")
	}
}

func TestWithCORSShouldAllowPreflightForUpsertMethod(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/users/42", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", http.MethodPut)
	recorder := httptest.NewRecorder()

	newCORSHandler().ServeHTTP(recorder, req)

	if got := recorder.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, http.MethodPut) {
		t.Fatalf("expected PUT in allow-methods, got %q", got)
	}
}

func TestWithCORSShouldNotLeakHeadersOnDeniedPreflight(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodOptions, "/users", nil)
	req.Header.Set("Origin", "http://evil.example")
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	recorder := httptest.NewRecorder()

	newCORSHandler().ServeHTTP(recorder, req)

	if got := recorder.Header().Get("Access-Control-Allow-Methods"); got != "" {
		t.Fatalf("expected no allow-methods for denied preflight, got %q", got)
	}
}

func TestWithCORSShouldVaryOnOriginWithoutOriginHeader(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	newCORSHandler().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected same-origin request to pass through, got %d", recorder.Code)
	}
	if got := recorder.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("expected Vary: Origin, got %q", got)
	}
}
