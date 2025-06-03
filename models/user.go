package models

type UserSignup struct {
	ID       int    `json:"id"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32,alphanum"`
	Username string `json:"username" validate:"required"`
	Message  string `json:"message"`
}

type UserLogin struct {
	ID       string `json:"id" db:"id" `
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
