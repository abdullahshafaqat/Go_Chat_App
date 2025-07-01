package db

import (
	"errors"

	"github.com/abdullahshafaqat/Go_Chat_App.git/models"
	"github.com/gin-gonic/gin"
)

func (d *dbImpl) CreateUser(c *gin.Context, user *models.UserSignup) error {
	var exist string
	err := d.db.QueryRow("SELECT email FROM users WHERE email = $1", user.Email).Scan(&exist)
	if err == nil {
		return errors.New("user already exists with this email")
	}
	query :=`INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = d.db.QueryRowxContext(c, query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
