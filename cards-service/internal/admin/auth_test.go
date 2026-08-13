package admin

import "testing"

func TestTokenMatchesShouldAcceptExactBearerToken(t *testing.T) {
	if !tokenMatches("secret", "Bearer secret") {
		t.Fatal("want exact token to pass")
	}
}

func TestTokenMatchesShouldRejectEmptyConfiguredToken(t *testing.T) {
	if tokenMatches("", "Bearer secret") {
		t.Fatal("want empty ADMIN_API_TOKEN to close the admin api")
	}
}

func TestTokenMatchesShouldRejectMissingHeader(t *testing.T) {
	if tokenMatches("secret", "") {
		t.Fatal("want missing Authorization header to be rejected")
	}
}

func TestTokenMatchesShouldRejectOtherScheme(t *testing.T) {
	if tokenMatches("secret", "Basic secret") {
		t.Fatal("want non bearer scheme to be rejected")
	}
}

func TestTokenMatchesShouldRejectDifferentToken(t *testing.T) {
	if tokenMatches("secret", "Bearer other-secret") {
		t.Fatal("want different token to be rejected")
	}
}
