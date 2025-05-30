package authservice

import (
	authservice "github.com/abdullahshafaqat/Go_ChatApp.git/api"
	db "github.com/abdullahshafaqat/Go_Chat_App.git/db/postgres"
)

type AuthService struct {
	userAuth db.DBStorage
}

type NewAuthServiceImpl struct {
	UserAuth db.DBStorage
}

func NewAuthService(input NewAuthServiceImpl) authservice.AuthService {
	return &AuthService{
		userAuth: input.UserAuth,
	}
}

var _ authservice.AuthService = &AuthService{}
