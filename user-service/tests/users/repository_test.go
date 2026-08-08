package users_test

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"user-service/internal/users"
)

func TestMapScanErrorShouldReturnNotFoundWhenNoRows(t *testing.T) {
	err := users.MapScanError(sql.ErrNoRows)

	if !errors.Is(err, users.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMapScanErrorShouldWrapUnexpectedErrors(t *testing.T) {
	original := errors.New("connection reset")

	err := users.MapScanError(original)

	if errors.Is(err, users.ErrNotFound) {
		t.Fatal("did not expect ErrNotFound for connection error")
	}
	if !errors.Is(err, original) {
		t.Fatalf("expected wrapped original error, got %v", err)
	}
	if fmt.Sprint(err) == original.Error() {
		t.Fatal("expected wrapped error message")
	}
}
