package models

type Playlist struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
	Tracks []int  `json:"tracks"` // ID треков
}
