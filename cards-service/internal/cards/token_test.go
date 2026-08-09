package cards

import "testing"

func TestShareTokenShouldRoundTrip(t *testing.T) {
	key := []byte("test-signing-key")

	token := GenerateShareToken(key, "42", 2024)

	id, year, err := DecodeShareToken(key, token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if id != "42" {
		t.Fatalf("expected user id 42, got %q", id)
	}
	if year != 2024 {
		t.Fatalf("expected year 2024, got %d", year)
	}
}

func TestShareTokenShouldPreserveUserIDWithColons(t *testing.T) {
	key := []byte("test-signing-key")

	token := GenerateShareToken(key, "tenant:eu:9", 2025)

	id, year, err := DecodeShareToken(key, token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if id != "tenant:eu:9" {
		t.Fatalf("expected user id tenant:eu:9, got %q", id)
	}
	if year != 2025 {
		t.Fatalf("expected year 2025, got %d", year)
	}
}

func TestShareTokenShouldRejectWrongKey(t *testing.T) {
	token := GenerateShareToken([]byte("real-key"), "42", 2024)

	if _, _, err := DecodeShareToken([]byte("other-key"), token); err == nil {
		t.Fatal("expected signature error for wrong key")
	}
}

func TestShareTokenShouldRejectTamperedPayload(t *testing.T) {
	key := []byte("test-signing-key")

	token := GenerateShareToken(key, "42", 2024)
	forged := GenerateShareToken(key, "999", 2024)

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
