package llm

import (
	"regexp"
	"strings"
)

var whitespaceRe = regexp.MustCompile(`\s+`)

var enclosingQuotes = []string{`"`, "'", "`", "“", "”", "«", "»"}

func normalizeText(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "- ")
	s = strings.ReplaceAll(s, "**", "")
	s = strings.ReplaceAll(s, "`", "")
	s = trimQuotes(s)
	s = whitespaceRe.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}

func trimQuotes(s string) string {
	s = strings.TrimSpace(s)
	for {
		before := s
		for _, q := range enclosingQuotes {
			s = strings.TrimPrefix(s, q)
			s = strings.TrimSuffix(s, q)
		}
		s = strings.TrimSpace(s)
		if s == before {
			return s
		}
	}
}

func stripCodeFence(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "```") {
		return s
	}
	if idx := strings.IndexByte(s, '\n'); idx != -1 {
		s = s[idx+1:]
	} else {
		s = strings.TrimPrefix(s, "```")
	}
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func extractJSONObject(s string) string {
	start := strings.IndexByte(s, '{')
	end := strings.LastIndexByte(s, '}')
	if start == -1 || end == -1 || end < start {
		return s
	}
	return s[start : end+1]
}
