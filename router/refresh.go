package router

import (
	"net/http"

	"github.com/abdullahshafaqat/Go_Chat_App.git/middelwares"
	"github.com/gin-gonic/gin"
)

func (r *routerImpl) RefreshToken(c *gin.Context) {
	refreshToken := c.PostForm("refresh_token")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	userID, err := middelwares.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	accessToken, _, err := middelwares.GenerateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
