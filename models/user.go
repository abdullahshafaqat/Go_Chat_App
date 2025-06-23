package models

type UserSignup struct {
    ID       string `json:"id"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=32,alphanum"`
    Username string `json:"username" validate:"required"`
}

type UserLogin struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}