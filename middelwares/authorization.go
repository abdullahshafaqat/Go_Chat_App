package middelwares

import (
	"fmt"
	"log"
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

// Update your AuthMiddleware to be more verbose
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

		// First try to parse as access token
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

		// Then check if it's a refresh token
		refreshToken, err := parseToken(tokenString, refreshKey)
		if err == nil && tokenType(refreshToken, "refresh") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":    "Refresh token not allowed for API access",
				"solution": "Use /refresh endpoint to get new access token",
			})
			return
		}

		// Handle specific error cases
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

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   userID,
		"type": "access",
		"exp":  time.Now().Add(time.Minute * 60).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(accessKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":   userID,
		"type": "refresh",
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}
func VerifyRefreshToken(tokenString string) (string, error) {
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

		log.Printf("WebSocket connection attempt from %s", c.Request.RemoteAddr)

		var tokenString string
		tokenString = BearerToken(c.GetHeader("Authorization"))

		if tokenString == "" {
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			log.Printf("No token provided for WebSocket connection from %s", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":    "Missing authorization token",
				"solution": "Include 'Authorization: Bearer <token>' header OR '?token=<token>' query parameter",
			})
			return
		}

		token, err := parseToken(tokenString, accessKey)
		if err != nil || !token.Valid {
			errorMsg := "Invalid token"
			if err != nil && strings.Contains(err.Error(), "expired") {
				errorMsg = "Token expired"
			}
			log.Printf("Invalid token from %s: %v", c.Request.RemoteAddr, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   errorMsg,
				"details": fmt.Sprintf("%v", err),
			})
			return
		}

		if !tokenType(token, "access") {
			log.Printf("Wrong token type from %s", c.Request.RemoteAddr)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":    "Invalid token type",
				"solution": "Use an access token for WebSocket connections",
			})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, exists := claims["ID"]; exists {

				userIDStr := fmt.Sprintf("%v", userID)
				c.Set("userID", userIDStr)
				log.Printf("Authenticated WebSocket connection for user %s", userIDStr)
				c.Next()
				return
			}
		}

		log.Printf("Invalid token claims from %s", c.Request.RemoteAddr)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})
	}
}

// func BackendWSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenString := BearerToken(c.GetHeader("Authorization"))
// 		if tokenString == "" {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error":    "Missing authorization header",
// 				"solution": "Include 'Authorization: Bearer <token>' header",
// 			})
// 			return
// 		}

// 		token, err := parseToken(tokenString, refreshKey)
// 		if err != nil || !token.Valid {
// 			errorMsg := "Invalid token"
// 			if err != nil && strings.Contains(err.Error(), "expired") {
// 				errorMsg = "Token expired"
// 			}
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error":   errorMsg,
// 				"details": err.Error(),
// 			})
// 			return
// 		}

// 		if !tokenType(token, "refresh") {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 				"error":    "Invalid token type",
// 				"solution": "Use a refresh token for backend WebSocket connections",
// 			})
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 			if userID, exists := claims["ID"]; exists {
// 				c.Set("userID", userID)
// 				c.Next()
// 				return
// 			}
// 		}

// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
// 			"error": "Invalid token claims",
// 		})
// 	}
// }
