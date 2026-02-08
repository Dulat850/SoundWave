package models

import "time"

type Artist struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	Bio        string    `json:"bio"`
	AvatarPath *string   `json:"avatar_path,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
