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

	accessToken, err := parseToken(tokenString, accessKey)
	if err == nil && TokenType(accessToken, "access") {
		return true, "access", nil
	} else if err != nil {
		log.Printf("Access token error: %v", err)
	}

	refreshToken, err := parseToken(tokenString, refreshKey)
	if err == nil {
		if TokenType(refreshToken, "refresh") {
			return false, "refresh token", nil
		}
	} else if strings.Contains(err.Error(), "expired") {
		return false, "", fmt.Errorf("token is expired")
	}

	return false, "", fmt.Errorf("invalid token")
}

func (s *serviceImpl) DecodeToken(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims format")
	}

	idClaim, exists := claims["ID"]
	if !exists {
		return "", fmt.Errorf("id claim not found in token")
	}

	id, ok := idClaim.(string)
	if !ok {
		return "", fmt.Errorf("id claim is not a string")
	}

	return id, nil
}
