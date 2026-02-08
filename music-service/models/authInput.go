package models

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Login    string `json:"login" binding:"required"` // сюда можно передать username ИЛИ email
	Password string `json:"password" binding:"required"`
}
