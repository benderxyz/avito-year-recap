package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"user-service/internal/users"
)

func TestParseUserIDShouldRejectZero(t *testing.T) {
	_, err := parseUserID("0")
	if err == nil {
		t.Fatal("expected error for zero user id")
	}
}

func TestParseUserIDShouldAcceptPositiveUint64(t *testing.T) {
	got, err := parseUserID("42")
	if err != nil {
		t.Fatalf("expected valid user id, got %v", err)
	}
	if got != 42 {
		t.Fatalf("expected user id 42, got %d", got)
	}
}

func TestDecodeJSONBodyShouldRejectUnknownFields(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/users/42", strings.NewReader(`{"username":"alex","unknown":true}`))
	recorder := httptest.NewRecorder()

	var body users.UpsertRequest
	err := decodeJSONBody(recorder, req, &body)

	if err == nil {
		t.Fatal("expected unknown field error")
	}
}

func TestDecodeJSONBodyShouldRejectMultipleObjects(t *testing.T) {
	req := httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/users/42", strings.NewReader(`{"username":"alex"} {"username":"bob"}`))
	recorder := httptest.NewRecorder()

	var body users.UpsertRequest
	err := decodeJSONBody(recorder, req, &body)

	if err == nil {
		t.Fatal("expected multiple object error")
	}
}
