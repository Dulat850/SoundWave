package models

type Artist struct {
	ID          int
	Name        string
	Tracks      []Track
	Albums      []Album
	Description any
}
