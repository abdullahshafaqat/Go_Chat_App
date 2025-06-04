package authservice

import (
	"errors"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *serviceImpl) Login(c *gin.Context, login *models.UserLogin) error {

	id, dbPassword, err := s.database.GetUserByEmail(login.Email)
	if err != nil {
		return errors.New("user not found")
	}


	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(login.Password)); err != nil {
		return errors.New("invalid password")
	}
	login.ID = id
	return nil
}
