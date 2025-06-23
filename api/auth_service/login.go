package authservice

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (s *serviceImpl) Login(c *gin.Context, email, password string) (string, error) {
	id, dbPassword, err := s.database.GetUserByEmail(email)
	if err != nil {
		log.Printf("Login failed for %s: user not found", email)
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		log.Printf("Login failed for %s: invalid password", email)
		return "", errors.New("invalid credentials")
	}

	log.Printf("Login successful for %s (ID: %s)", email, id)
	return id, nil
}
