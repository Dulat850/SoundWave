package models

type Album struct {
	ID       int
	Title    string
	ArtistID int
	Tracks   []Track
	CoverURL any
}
