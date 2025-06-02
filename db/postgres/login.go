package db

import (
	"log"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (d *dbImpl) GetUserByEmail(c *gin.Context, email string) (*models.UserLogin, error) {
	var user models.UserLogin

	err := d.db.QueryRow("SELECT id, email, password FROM signup WHERE email=$1", email).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		log.Fatal("Error: User not exist")
		return nil, err
	}

	return &user, nil

}
