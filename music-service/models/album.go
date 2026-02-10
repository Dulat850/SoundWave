package models

import "time"

type Album struct {
	ID         int       `json:"id"`
	ArtistID   int       `json:"artist_id"`
	Title      string    `json:"title"`
	CoverPath  *string   `json:"cover_path,omitempty"`
	ReleasedAt *string   `json:"released_at,omitempty"` // YYYY-MM-DD (можно потом сделать time.Time)
	CreatedAt  time.Time `json:"created_at"`
	CoverURL   *string   `json:"cover_url,omitempty"`
}
