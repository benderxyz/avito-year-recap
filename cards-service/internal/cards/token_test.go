package cards

import "testing"

func TestShareTokenShouldRoundTrip(t *testing.T) {
	key := []byte("test-signing-key")

	token := GenerateShareToken(key, "user-123", 2024)

	id, year, err := DecodeShareToken(key, token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if id != "user-123" {
		t.Fatalf("expected external id user-123, got %q", id)
	}
	if year != 2024 {
		t.Fatalf("expected year 2024, got %d", year)
	}
}

func TestShareTokenShouldPreserveExternalIDWithColons(t *testing.T) {
	key := []byte("test-signing-key")

	token := GenerateShareToken(key, "tenant:eu:user-9", 2025)

	id, year, err := DecodeShareToken(key, token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if id != "tenant:eu:user-9" {
		t.Fatalf("expected external id tenant:eu:user-9, got %q", id)
	}
	if year != 2025 {
		t.Fatalf("expected year 2025, got %d", year)
	}
}

func TestShareTokenShouldRejectWrongKey(t *testing.T) {
	token := GenerateShareToken([]byte("real-key"), "user-123", 2024)

	if _, _, err := DecodeShareToken([]byte("other-key"), token); err == nil {
		t.Fatal("expected signature error for wrong key")
	}
}

func TestShareTokenShouldRejectTamperedPayload(t *testing.T) {
	key := []byte("test-signing-key")

	token := GenerateShareToken(key, "user-123", 2024)
	forged := GenerateShareToken(key, "user-999", 2024)

	_, mac, _ := lastCut(token, shareTokenSeparator)
	tampered, _, _ := lastCut(forged, shareTokenSeparator)

	if _, _, err := DecodeShareToken(key, tampered+shareTokenSeparator+mac); err == nil {
		t.Fatal("expected signature error for tampered payload")
	}
}

func TestShareTokenShouldRejectMalformedToken(t *testing.T) {
	if _, _, err := DecodeShareToken([]byte("k"), "no-separator-here"); err == nil {
		t.Fatal("expected error for token without separator")
	}
}
