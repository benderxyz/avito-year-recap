package users

import "time"

type User struct {
	UserID     uint64    `json:"user_id"`
	ExternalID string    `json:"external_id"`
	Username   string    `json:"username"`
	Timezone   string    `json:"timezone"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpsertRequest struct {
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Timezone   string `json:"timezone"`
}
