package llm

import (
	"strings"
	"unicode"
)

const (
	maxBadgeTitleLen   = 40
	maxBadgeDescLen    = 120
	maxInsightTitleLen = 60
	maxInsightTextLen  = 200
	maxDisplayNameLen  = 50
)

func validateShareableText(s string, maxLen int) (string, bool) {
	if s == "" {
		return "", false
	}
	if runeLen(s) > maxLen {
		return "", false
	}
	if strings.ContainsAny(s, "\n\r") {
		return "", false
	}
	if containsDigit(s) {
		return "", false
	}

	return s, true
}

func validateBadge(t BadgeText) (BadgeText, bool) {
	title, okTitle := validateShareableText(normalizeText(t.Title), maxBadgeTitleLen)
	desc, okDesc := validateShareableText(normalizeText(t.Description), maxBadgeDescLen)
	if !okTitle || !okDesc {
		return BadgeText{}, false
	}
	return BadgeText{Title: title, Description: desc}, true
}

func validateInsight(in InsightText) (InsightText, bool) {
	title, okTitle := validateShareableText(normalizeText(in.Title), maxInsightTitleLen)
	text, okText := validateShareableText(normalizeText(in.Text), maxInsightTextLen)
	if !okTitle || !okText {
		return InsightText{}, false
	}
	return InsightText{Title: title, Text: text}, true
}

func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func runeLen(s string) int {
	return len([]rune(s))
}
