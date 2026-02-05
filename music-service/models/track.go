package models

type Track struct {
	ID       int
	Title    string
	Duration int
	category string
	AlbumID  int
	AudioURL any
	ArtistID any
}
