package models

import "time"

type Track struct {
	ID              int       `json:"id"`
	ArtistID        int       `json:"artist_id"`
	AlbumID         *int      `json:"album_id,omitempty"`
	GenreID         *int      `json:"genre_id,omitempty"`
	Title           string    `json:"title"`
	DurationSeconds int       `json:"duration_seconds"`
	AudioPath       string    `json:"audio_path"`
	CoverPath       *string   `json:"cover_path,omitempty"`
	PlayCount       int64     `json:"play_count"`
	CreatedAt       time.Time `json:"created_at"`
	Duration        string    `json:"duration,omitempty"` // "mm:ss"
	AudioURL        *string   `json:"audio_url,omitempty"`
	CoverURL        *string   `json:"cover_url,omitempty"`
}
