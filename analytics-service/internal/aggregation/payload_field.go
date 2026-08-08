package aggregation

import "regexp"

var payloadFieldPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func sanitizePayloadField(field string) string {
	if field == "" {
		field = "value"
	}
	if !payloadFieldPattern.MatchString(field) {
		return "value"
	}
	return field
}
