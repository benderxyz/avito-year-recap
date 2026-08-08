package models

type RecapMode string

const (
	RecapModePrivate RecapMode = "private"
	RecapModePublic  RecapMode = "public"
)

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
	Features      *Features              `json:"features,omitempty"`
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

type Profile struct {
	UserID     uint64 `json:"user_id"`
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Timezone   string `json:"timezone"`
}

type MetricValue struct {
	Type     MetricType `json:"type"`
	Value    any        `json:"value"`
	Currency Currency   `json:"currency,omitempty"`
}

type Badge struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl,omitempty"`
}

type Features struct {
	ShareEnabled bool   `json:"shareEnabled"`
	ShareURL     string `json:"shareUrl,omitempty"`
	UpsellID     string `json:"upsellId,omitempty"`
}
