package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`    // не в JSON
	Role     string `json:"role"` // "user" или "admin"
}
