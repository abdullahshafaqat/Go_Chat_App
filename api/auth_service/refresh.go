package authservice

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *serviceImpl) GenerateTokens(ID string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   ID,
		"type": "access",
		"exp":  time.Now().Add(time.Minute * 1).Unix(),
	})
	fmt.Println("Generating tokens for ID:", ID)
	accessTokenString, err := accessToken.SignedString(accessKey)
	if err != nil {
		return "", "", err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   ID,
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil

}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return accessKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
func (s *serviceImpl) VerifyRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return refreshKey, nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["ID"] == nil {
		return "", fmt.Errorf("invalid claims in token")
	}

	id := claims["ID"].(string)
	return id, nil

}
