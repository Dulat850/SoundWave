package models

type Album struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	ArtistID int    `json:"artist_id"`
	CoverURL string `json:"cover_url"`
}

//
