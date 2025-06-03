package authservice

import (
	"fmt"
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

	if token, err := parseToken(tokenString, accessKey); err == nil {
		if TokenType(token, "access") {
			return true, "access", nil
		}
	}

	if token, err := parseToken(tokenString, refreshKey); err == nil {
		if TokenType(token, "refresh") {
			return false, "refresh token ", nil
		}
	} else if strings.Contains(err.Error(), "expired") {
		return false, "", fmt.Errorf("token is expired")
	}

	return false, "", fmt.Errorf("invalid token")
}
