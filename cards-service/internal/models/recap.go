package models

type Badge struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Recap struct {
	ProfileID string  `json:"profile_id"`
	Badges    []Badge `json:"badges"`
}
