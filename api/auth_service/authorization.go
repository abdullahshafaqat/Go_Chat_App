package authservice

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var accessKey = []byte(os.Getenv("AC_SECRET"))
var refreshKey = []byte(os.Getenv("RF_SECRET"))

func (s *serviceImpl) BearerToken(header string) string {
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}

func parseToken(tokenString string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}

func TokenType(token *jwt.Token, Type string) bool {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, ok := claims["type"].(string); ok && tokenType == Type {
			return true
		}
	}
	return false
}
func (s *serviceImpl) Authorize(token string) (bool, string, error) {
	tokenString := token
	if tokenString == "" {
		return false, "", fmt.Errorf("missing or invalid authorization header")
	}

	// Check Access Token
	accessToken, err := parseToken(tokenString, accessKey)
	if err == nil && TokenType(accessToken, "access") {
		return true, "access", nil
	} else if err != nil {
		log.Printf("Access token error: %v", err) // Debug log
	}

	// Check Refresh Token
	refreshToken, err := parseToken(tokenString, refreshKey)
	if err == nil && TokenType(refreshToken, "refresh") {
		return false, "refresh token", nil
	} else if err != nil {
		log.Printf("Refresh token error: %v", err) // Debug log
	}

	return false, "", fmt.Errorf("invalid token")
}
