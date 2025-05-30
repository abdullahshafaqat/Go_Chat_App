package authservice

import (
	"fmt"

	"github.com/abdullahshafaqat/Go_ChatApp.git/models"
	"github.com/gin-gonic/gin"
)

func (s *AuthService) SignUp(c *gin.Context, req *models.UserSignup) *models.UserSignup {

	User := s.userAuth.SignUp(c, req)

	if User == nil {
		fmt.Println("signup error")
		return nil
	}

	res := models.UserSignup{
		Email:    User.Email,
		Password: User.Password,
		Message:  "User created successfully",
	}

	return &res
}
