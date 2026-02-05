package models

import "time"

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string `json:"role"`
}

//
