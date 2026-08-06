package models

type Recap struct {
	ProfileID string  `json:"profile_id"`
	Badges    []Badge `json:"badges"`
}

type RecapPayload struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Meta          Meta                   `json:"meta"`
	Metrics       map[string]MetricValue `json:"metrics"`
	Badges        []Badge                `json:"badges"`
	Story         []map[string]any       `json:"story"` // пока просто map, типы сцен добавить позже
}

type Meta struct {
	Vertical    string `json:"vertical"`
	Year        int    `json:"year"`
	Locale      string `json:"locale"`
	User        User   `json:"user"`
	GeneratedAt string `json:"generatedAt"`
}

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
}

type MetricValue struct {
	Type     string `json:"type"`               // "number" | "money" | "percentile" | ...
	Value    any    `json:"value"`              // число или строка
	Currency string `json:"currency,omitempty"` // только для money
}

type Badge struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl,omitempty"`
}
