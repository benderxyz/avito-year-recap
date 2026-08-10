package users

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const DefaultTimezone = "UTC"

var ErrInvalidTimezone = errors.New("invalid timezone")

func NormalizeTimezone(raw string) (string, error) {
	timezone := strings.TrimSpace(raw)
	if timezone == "" {
		return DefaultTimezone, nil
	}

	if timezone == "Local" {
		return "", fmt.Errorf("%w: %q depends on the host and is not portable", ErrInvalidTimezone, timezone)
	}

	if _, err := time.LoadLocation(timezone); err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidTimezone, timezone)
	}

	return timezone, nil
}
