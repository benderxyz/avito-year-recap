package models

type User struct {
	UserID     uint64 `json:"user_id"`
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Timezone   string `json:"timezone"`
}
