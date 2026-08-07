package cards

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const shareTokenSeparator = "."

func GenerateShareToken(signingKey []byte, externalID string, year int) string {
	payload := shareTokenPayload(externalID, year)
	mac := signPayload(signingKey, payload)

	return strings.Join([]string{
		base64.RawURLEncoding.EncodeToString([]byte(payload)),
		base64.RawURLEncoding.EncodeToString(mac),
	}, shareTokenSeparator)
}

func DecodeShareToken(signingKey []byte, token string) (externalID string, year int, err error) {
	encodedPayload, encodedMAC, found := strings.Cut(token, shareTokenSeparator)
	if !found {
		return "", 0, fmt.Errorf("malformed share token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return "", 0, fmt.Errorf("decode payload: %w", err)
	}

	gotMAC, err := base64.RawURLEncoding.DecodeString(encodedMAC)
	if err != nil {
		return "", 0, fmt.Errorf("decode signature: %w", err)
	}

	wantMAC := signPayload(signingKey, string(payload))
	if !hmac.Equal(gotMAC, wantMAC) {
		return "", 0, fmt.Errorf("invalid share token signature")
	}

	id, rawYear, found := lastCut(string(payload), ":")
	if !found {
		return "", 0, fmt.Errorf("malformed share token payload")
	}

	year, err = strconv.Atoi(rawYear)
	if err != nil {
		return "", 0, fmt.Errorf("parse year: %w", err)
	}

	return id, year, nil
}

func shareTokenPayload(externalID string, year int) string {
	return externalID + ":" + strconv.Itoa(year)
}

func signPayload(signingKey []byte, payload string) []byte {
	mac := hmac.New(sha256.New, signingKey)
	mac.Write([]byte(payload))
	return mac.Sum(nil)
}

func lastCut(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
