package models

type Favorite struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	TrackID int `json:"track_id"`
}
