package models

import "time"

type Listen struct {
	ID         int       `json:"id"`
	UserID     *int      `json:"user_id,omitempty"`
	TrackID    int       `json:"track_id"`
	ListenedAt time.Time `json:"listened_at"`
}
