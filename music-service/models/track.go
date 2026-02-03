package models

type Track struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	ArtistID int    `json:"artist_id"`
	AlbumID  int    `json:"album_id"`
	Duration int    `json:"duration"` // секунды
	AudioURL string `json:"audio_url"`
}

//
