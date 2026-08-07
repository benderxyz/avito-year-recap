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
	Story         []map[string]any       `json:"story"`
}

type Meta struct {
	Vertical    string  `json:"vertical"`
	Year        int     `json:"year"`
	Locale      string  `json:"locale"`
	Profile     Profile `json:"profile"`
	GeneratedAt string  `json:"generatedAt"`
}

type Profile struct {
	UserID     uint64 `json:"user_id"`
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Timezone   string `json:"timezone"`
}

type MetricValue struct {
	Type     string `json:"type"`
	Value    any    `json:"value"`
	Currency string `json:"currency,omitempty"`
}

type Badge struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl,omitempty"`
}
