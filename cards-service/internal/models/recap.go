package models

type Card struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Text        string `json:"text"`
	Action      string `json:"action,omitempty"`
	ActionValue string `json:"action_value,omitempty"`
}

type Recap struct {
	ProfileID string `json:"profile_id"`
	Cards     []Card `json:"cards"`
}
