package db

import (
	"net/http"

	"github.com/abdullahshafaqat/Go_ChatApp.git/models"
	"github.com/gin-gonic/gin"
)

func (d *StorageImpl) SignUp(c *gin.Context, req *models.UserSignup) *models.UserSignup {

	err := d.db.QueryRow("SELECT email FROM users WHERE email = $1", &req.Email).Scan(&req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists with this email"})
		return nil
	}
	_, err = d.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", &req.Email, &req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return nil
	}

	return req

}
