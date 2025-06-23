package middelwares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	accessKey  = []byte(os.Getenv("AC_SECRET"))
	refreshKey = []byte(os.Getenv("RF_SECRET"))
)

func BearerToken(header string) string {
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

func tokenType(token *jwt.Token, expectedType string) bool {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, ok := claims["type"].(string); ok && tokenType == expectedType {
			return true
		}
	}
	return false
}


func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := BearerToken(c.GetHeader("Authorization"))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":    "Missing authorization header",
				"solution": "Include 'Authorization: Bearer <token>' header",
			})
			return
		}

		
		accessToken, err := parseToken(tokenString, accessKey)
		if err == nil && tokenType(accessToken, "access") {
			if claims, ok := accessToken.Claims.(jwt.MapClaims); ok {
				if userID, exists := claims["ID"]; exists {
					c.Set("userID", userID)
					c.Next()
					return
				}
			}
		}

	
		refreshToken, err := parseToken(tokenString, refreshKey)
		if err == nil && tokenType(refreshToken, "refresh") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":    "Refresh token not allowed for API access",
				"solution": "Use /refresh endpoint to get new access token",
			})
			return
		}

		
		if err != nil {
			errorMsg := "Invalid token"
			if strings.Contains(err.Error(), "expired") {
				errorMsg = "Token expired"
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   errorMsg,
				"details": err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token",
		})
	}
}

func GenerateTokens(userID string) (string, string, error) {
	now := time.Now()


	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":            userID,
		"type":          "access",
		"iat":           now.Unix(),
		"exp":           now.Add(50 * time.Minute).Unix(),
		"token_version": 1,
	})

	accessTokenString, err := accessToken.SignedString(accessKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":            userID,
		"type":          "refresh",
		"iat":           now.Unix(),
		"exp":           now.Add(7 * 24 * time.Hour).Unix(),
		"token_version": 1,
	})

	refreshTokenString, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}
func VerifyRefreshToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("token validation failed: %v", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	if claims["ID"] == nil || claims["type"] != "refresh" {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims["ID"].(string), nil
}

func getClaimID(token *jwt.Token) string {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if id, ok := claims["ID"].(string); ok {
			return id
		}
	}
	return ""
}

func GetUserIDFromToken(tokenString string) (string, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	return getClaimID(token), nil
}

func WSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if tokenString == "" {
			tokenString = c.Query("token")
		}

	
		if tokenString == "" || tokenString == "undefined" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			return
		}

		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return accessKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, exists := claims["ID"].(string)
		if !exists || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		// 4. Set user ID in context
		c.Set("userID", userID)
		c.Next()
	}
}
