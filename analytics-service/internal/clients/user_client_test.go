package clients_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"analytics-service/internal/clients"
)

func TestTimezoneShouldReturnUserTimezoneWhenUserExists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/users/42" {
			t.Fatalf("expected path /users/42, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user_id":42,"external_id":"x","username":"Alex","timezone":"Europe/Moscow"}`))
	}))
	defer server.Close()

	client := clients.NewUserClient(server.URL)

	timezone, err := client.Timezone(context.Background(), 42)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if timezone != "Europe/Moscow" {
		t.Fatalf("expected Europe/Moscow, got %s", timezone)
	}
}

func TestTimezoneShouldReturnUTCWhenUserNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"user not found"}`))
	}))
	defer server.Close()

	client := clients.NewUserClient(server.URL)

	timezone, err := client.Timezone(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if timezone != "UTC" {
		t.Fatalf("expected UTC, got %s", timezone)
	}
}

func TestTimezoneShouldReturnUTCWhenTimezoneIsEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"user_id":1,"external_id":"","username":"","timezone":"  "}`))
	}))
	defer server.Close()

	client := clients.NewUserClient(server.URL)

	timezone, err := client.Timezone(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if timezone != "UTC" {
		t.Fatalf("expected UTC, got %s", timezone)
	}
}

func TestGetUserShouldReturnNotFoundWhenStatusIs404(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := clients.NewUserClient(server.URL)

	_, err := client.GetUser(context.Background(), 9)
	if !errors.Is(err, clients.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestTimezoneShouldFailWhenUserServiceReturnsErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := clients.NewUserClient(server.URL)

	_, err := client.Timezone(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got != "resolve timezone: user-service returned status 500" {
		t.Fatalf("unexpected error message: %s", got)
	}
}
