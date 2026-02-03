package models

type Playlist struct {
	ID     int
	Name   string
	UserID int
	Tracks []Track
}
