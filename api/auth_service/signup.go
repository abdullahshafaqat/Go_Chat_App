package authservice

import (
	"errors"
	"regexp"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *serviceImpl) SignUp(c *gin.Context, Newuser *models.UserSignup) error {
	if !isGmail(Newuser.Email) {
		return errors.New("please enter valid email address")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Newuser.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	Newuser.Password = string(hashedPassword)
	return s.db.CreateUser(c, Newuser)
}

func isGmail(email string) bool {
	return regexp.MustCompile(`^[^@]+@gmail\.com$`).MatchString(email)
}
