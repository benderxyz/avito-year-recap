package users_test

import (
	"errors"
	"testing"

	"user-service/internal/users"
)

func TestNormalizeTimezoneShouldDefaultWhenEmpty(t *testing.T) {
	for _, raw := range []string{"", "   ", "\t"} {
		got, err := users.NormalizeTimezone(raw)
		if err != nil {
			t.Fatalf("expected no error for %q, got %v", raw, err)
		}
		if got != users.DefaultTimezone {
			t.Fatalf("expected %q for %q, got %q", users.DefaultTimezone, raw, got)
		}
	}
}

func TestNormalizeTimezoneShouldAcceptKnownZones(t *testing.T) {
	for _, raw := range []string{"UTC", "Europe/Moscow", "  Europe/Moscow  ", "America/New_York"} {
		got, err := users.NormalizeTimezone(raw)
		if err != nil {
			t.Fatalf("expected %q to be accepted, got %v", raw, err)
		}
		if got == "" {
			t.Fatalf("expected normalized value for %q", raw)
		}
	}
}

func TestNormalizeTimezoneShouldRejectUnknownZones(t *testing.T) {
	for _, raw := range []string{"Mars/Olympus", "Europe/Moskva", "'; DROP TABLE users;--", "+03:00"} {
		got, err := users.NormalizeTimezone(raw)
		if !errors.Is(err, users.ErrInvalidTimezone) {
			t.Fatalf("expected ErrInvalidTimezone for %q, got %v", raw, err)
		}
		if got != "" {
			t.Fatalf("expected empty result for rejected %q, got %q", raw, got)
		}
	}
}

func TestNormalizeTimezoneShouldRejectHostDependentLocal(t *testing.T) {
	if _, err := users.NormalizeTimezone("Local"); !errors.Is(err, users.ErrInvalidTimezone) {
		t.Fatalf("expected Local to be rejected, got %v", err)
	}
}
