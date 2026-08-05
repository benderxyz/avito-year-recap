package users

import "time"

func ApplyUpsertTimestamps(user User, existing *User, now time.Time) User {
	if existing != nil {
		user.CreatedAt = existing.CreatedAt
	} else if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}

	user.UpdatedAt = now
	if user.Timezone == "" {
		user.Timezone = "UTC"
	}

	return user
}
